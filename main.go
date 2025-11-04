package main

import (
	"flag"
	"fmt"
	"github.com/ZeroHawkeye/wordZero/pkg/document"
	"github.com/ZeroHawkeye/wordZero/pkg/style"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Dir                string   // 读取的代码文件目录
	IncludeFileSuffix  []string // 需要包含的文件后缀
	IncludeFileSuffixs string
	Name               string // 页眉上显示的名字
}

var defaultConfig = Config{}

var fileList []*FileItem

func init() {
	flag.StringVar(&defaultConfig.Dir, "r", "", "程序目录")
	flag.StringVar(&defaultConfig.IncludeFileSuffixs, "f", ".go", "需要包含的文件，多个如下 .go;.vue")
	flag.StringVar(&defaultConfig.Name, "n", "", "页眉的名称")
	flag.Parse()
}

func main() {
	// 示例 ./main -r D:\IT\fengfeng\test -f .go -n "test v1.0"
	FlagsRun()
	// 读取文件列表
	err := WalkDir()
	if err != nil {
		fmt.Println(err)
		return
	}
	DocHandler()
}

func FlagsRun() {
	if defaultConfig.Dir == "" {
		fmt.Println("请输入目录")
		os.Exit(0)
	}
	defaultConfig.IncludeFileSuffix = strings.Split(defaultConfig.IncludeFileSuffixs, ";")
}

func DocHandler() {
	// 创建word文档
	doc := document.New()

	// 页眉+页脚
	err := doc.AddHeader(document.HeaderFooterTypeDefault, defaultConfig.Name)
	if err != nil {
		log.Printf("添加页眉失败: %v", err)
		return
	}
	err = doc.AddFooterWithPageNumber(document.HeaderFooterTypeDefault, "", true)
	if err != nil {
		log.Printf("添加带页码的页脚失败: %v", err)
		return
	}

	chapter1 := doc.AddParagraph("前30页")
	chapter1.SetStyle(style.StyleHeading2)
	for _, item := range fileList {
		fmt.Println("读取文件", item.EPath)
		byteData, err := os.ReadFile(item.Path)
		if err != nil {
			fmt.Println("读取文件失败", err)
			return
		}
		lines := strings.Split(string(byteData), "\n")

		text := doc.AddParagraph(item.EPath)
		text.SetStyle(style.StyleNormal)
		for _, line := range lines {
			para := doc.AddParagraph("")
			para.AddFormattedText(line, &document.TextFormat{
				FontFamily: "等线 (西文正文)",
				FontSize:   11,
			})
		}
		paragraphs := doc.Body.GetParagraphs()
		page := int(math.Ceil(float64(len(paragraphs)) / float64(27)))
		if page >= 30 {
			fmt.Printf("当前页数:%d 超过30页 退出", page)
			break
		}
	}

	chapter2 := doc.AddParagraph("后30页")
	chapter2.SetStyle(style.StyleHeading2)
	ReverseSlice(fileList)
	for _, item := range fileList {
		fmt.Println("读取文件", item.EPath)
		byteData, err := os.ReadFile(item.Path)
		if err != nil {
			fmt.Println("读取文件失败", err)
			return
		}
		lines := strings.Split(string(byteData), "\n")

		text := doc.AddParagraph(item.EPath)
		text.SetStyle(style.StyleNormal)
		for _, line := range lines {
			para := doc.AddParagraph("")
			para.AddFormattedText(line, &document.TextFormat{
				FontFamily: "等线 (西文正文)",
				FontSize:   11,
			})
		}
		paragraphs := doc.Body.GetParagraphs()
		page := int(math.Ceil(float64(len(paragraphs)) / float64(27)))
		if page >= 60 {
			fmt.Printf("当前页数:%d 超过30页 退出", page)
			break
		}
	}

	doc.Save("程序鉴别材料.docx")
}

type FileItem struct {
	Path    string // 文件地址
	EPath   string // 文件的相对路径
	Content string // 文件的内容
}

func WalkDir() error {
	// 检查根目录是否存在
	info, err := os.Stat(defaultConfig.Dir)
	if err != nil {
		return fmt.Errorf("目录不存在或无法访问：%w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("不是有效的目录：%s", defaultConfig.Dir)
	}

	// 开始递归遍历（初始深度为0）
	return walkDirRecursive(defaultConfig.Dir, 0, defaultConfig)
}

func walkDirRecursive(dir string, currentDepth int, config Config) error {
	// 读取目录内容
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("读取目录[%s]失败：%w", dir, err)
	}

	// 遍历目录项
	for _, entry := range entries {
		// 跳过隐藏文件
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		entryPath := filepath.Join(dir, entry.Name())
		if !entry.IsDir() {
			// 如果这个后缀在列表里面
			var noContinue bool
			for _, suffix := range config.IncludeFileSuffix {
				if strings.HasSuffix(entry.Name(), suffix) {
					noContinue = true
				}
			}
			if !noContinue {
				continue
			}
			targetPath := strings.Replace(entryPath, config.Dir, "", 1)
			targetPath = strings.ReplaceAll(targetPath, "\\", "/")[1:]
			fileList = append(fileList, &FileItem{
				Path:  entryPath,
				EPath: targetPath,
			})
			continue
		}
		// 如果是目录，递归遍历
		walkDirRecursive(entryPath, currentDepth+1, config)
	}
	return nil
}

// ReverseSlice 原地反转切片（支持 []int、[]string、[]byte 等任意类型）
func ReverseSlice[T any](s []T) {
	left, right := 0, len(s)-1
	for left < right {
		// 交换左右指针指向的元素
		s[left], s[right] = s[right], s[left]
		// 指针向中间移动
		left++
		right--
	}
}
