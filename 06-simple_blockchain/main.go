package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

var BlockchainInstance *Blockchain

// Block 区块不需要有 json 的 tag 因为不用与 api 交互
type Block struct {
	Position  int           // 位置
	Data      *BookCheckout // 数据
	TimeStamp string        // 时间戳
	Hash      string        // 当前块哈希
	PrevHash  string        // 前一个区块哈希
}

type Blockchain struct {
	blocks []*Block
}

// Book 以下数据需要与 api 交互所以需要设置 json tag
type Book struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`       // 作者
	PublishDate string `json:"publish_date"` // 出版日期
	ISBN        string `json:"isbn"`         // 图书编码
}

// BookCheckout 图书结账
type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"` // 是否是创世块
}

// 创建 hash
func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data)

	data := strconv.Itoa(b.Position) + b.TimeStamp + string(bytes) + b.PrevHash

	hash := sha256.New()
	hash.Write([]byte(data))
	// 将一个字节片段哈希值转换为十六进制字符串
	// hash.Sum() 用于计算输入数据的哈希值，并返回一个字节片段
	b.Hash = hex.EncodeToString(hash.Sum(nil))
}

func CreateBlock(prevBlock *Block, checkoutItem *BookCheckout) *Block {
	block := new(Block)
	block.Position = prevBlock.Position + 1
	block.TimeStamp = time.Now().String()
	block.Data = checkoutItem
	block.PrevHash = prevBlock.Hash

	block.generateHash()

	return block
}

func validBlock(block, prevBlock *Block) bool {
	if block.PrevHash != prevBlock.Hash {
		return false
	}

	if block.Position != prevBlock.Position+1 {
		return false
	}

	if !block.validHash(block.Hash) {
		return false
	}
	return true
}

func (b *Block) validHash(hash string) bool {
	b.generateHash()
	if b.Hash != hash {
		return false
	}
	return true
}

func (bc *Blockchain) AddBlock(data *BookCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1]

	block := CreateBlock(prevBlock, data)

	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

func writeBlock(w http.ResponseWriter, r *http.Request) {
	checkoutItem := new(BookCheckout)
	// var checkoutItem BookCheckout
	if err := json.NewDecoder(r.Body).Decode(checkoutItem); err != nil {
		log.Printf("could not write block: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("could not write block"))
	}

	BlockchainInstance.AddBlock(checkoutItem)
	resp, err := json.MarshalIndent(checkoutItem, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not marshal payload: %v", err)
		_, _ = w.Write([]byte("could not write block"))
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(resp)
}

func newBook(w http.ResponseWriter, r *http.Request) {
	// var book Book
	book := new(Book)

	// 将 json 数据解码到 book 中
	if err := json.NewDecoder(r.Body).Decode(book); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("could not create book: %v\n", err)
		_, _ = w.Write([]byte("Create failed"))
		return
	}

	h := md5.New()
	// 向 h 中写入 book.ISBN + book.PublishDate 的字符串
	_, _ = io.WriteString(h, book.ISBN+book.PublishDate)
	// 通过 h.Sum(nil) 获取最终的哈希值 然后将其转换为十六进制字符串并赋值给 book.ID
	book.ID = fmt.Sprintf("%x", h.Sum(nil))

	// 将数据结构转换为 JSON 格式的字节数组
	res, err := json.MarshalIndent(book, "", " ")
	if err != nil {
		log.Printf("could not save book data: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(res)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Write response failed: %v\n", err)
		return
	}
}

func GenesisBlock() *Block {
	// block := new(Block)
	// checkout := new(BookCheckout)
	// checkout.IsGenesis = true
	// b := CreateBlock(block, checkout)
	// return b

	// 上下两种写法一样
	// 注意:
	return CreateBlock(&Block{}, &BookCheckout{IsGenesis: true})

}

func NewBlockchain() *Blockchain {
	// blockchain := new(Blockchain)
	// blockchain.blocks = append(blockchain.blocks, GenesisBlock())
	// return blockchain

	// 上面与下面的方式都可以

	return &Blockchain{[]*Block{GenesisBlock()}}
}

func getBlockchain(w http.ResponseWriter, r *http.Request) {
	jbyte, err := json.MarshalIndent(BlockchainInstance.blocks, "", " ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	_, _ = io.WriteString(w, string(jbyte))
}

func main() {
	BlockchainInstance = NewBlockchain()

	r := mux.NewRouter()
	r.HandleFunc("/", getBlockchain).Methods("GET")
	r.HandleFunc("/", writeBlock).Methods("POST")
	r.HandleFunc("/new", newBook).Methods("POST")

	go func() {
		for _, v := range BlockchainInstance.blocks {
			fmt.Printf("prev. hash:%x\n", v.PrevHash)
			fmt.Printf("Hash:%x\n", v.Hash)

			bytes, _ := json.MarshalIndent(v.Data, "", " ")
			fmt.Printf("Data:%v\n", string(bytes))
			fmt.Println()
		}
	}()

	log.Println("Listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", r))
}
