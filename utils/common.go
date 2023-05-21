package utils

import (
	"bufio"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"math/rand"
	"net"
	"net/mail"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"
)

// PrintError 错误输出
func PrintError(msg string) {
	fmt.Printf("\033[1;31m" + msg + " \033[0m\n")
}

// PrintSuccess 正确输出
func PrintSuccess(msg string) {
	fmt.Printf("\033[1;32m" + msg + " \033[0m\n")
}

// FormatYmdHis 返回格式：2021-08-05 00:00:01
func FormatYmdHis(timeObj time.Time) string {
	year := timeObj.Year()
	month := timeObj.Month()
	day := timeObj.Day()
	hour := timeObj.Hour()
	minute := timeObj.Minute()
	second := timeObj.Second()
	//注意：%02d 中的 2 表示宽度，如果整数不够 2 列就补上 0
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d", year, month, day, hour, minute, second)
}

// Exists 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// IsFile 判断所给路径是否为文件
func IsFile(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return !s.IsDir()
}

// Mkdir 创建目录
func Mkdir(path string, perm os.FileMode) (err error) {
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, perm)
		if err != nil {
			return
		}
		err = os.Chmod(path, perm)
		if err != nil {
			return
		}
	}
	return err
}

// ReadFile 读取文件
func ReadFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

// WriteFile 保存文件（string）
func WriteFile(path string, content string) error {
	return WriteByte(path, []byte(content))
}

// WriteByte 保存文件（byte）
func WriteByte(path string, fileByte []byte) error {
	fileDir := filepath.Dir(path)
	if !Exists(fileDir) {
		err := os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	return os.WriteFile(path, fileByte, 0666)
}

// AppendToFile 追加文件
func AppendToFile(path string, content string) error {
	fileDir := filepath.Dir(path)
	if !Exists(fileDir) {
		err := os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			return err
		}
	}
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if _, err = file.WriteString(content); err != nil {
		return err
	}
	return nil
}

// SliceInsert 向数组插入内容
func SliceInsert(s []string, index int, value string) []string {
	rear := append([]string{}, s[index:]...)
	return append(append(s[:index], value), rear...)
}

// FindIndex 查找数组位置
func FindIndex(tab []string, value string) int {
	for i, v := range tab {
		if v == value {
			return i
		}
	}
	return -1
}

// RandString 生成随机字符串
func RandString(len int) string {
	var r *rand.Rand
	r = rand.New(rand.NewSource(time.Now().Unix()))
	bs := make([]byte, len)
	for i := 0; i < len; i++ {
		b := r.Intn(26) + 65
		bs[i] = byte(b)
	}
	return string(bs)
}

// GenerateString 生成随机字符串
func GenerateString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var bs = make([]byte, length)
	_, err := rand.Read(bs)
	if err != nil {
		return RandString(length)
	}
	for i, b := range bs {
		bs[i] = charset[b%byte(len(charset))]
	}
	return string(bs)
}

// RandNum 生成随机数
func RandNum(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

// StringMd5 MD5
func StringMd5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// StringMd52 MD5
func StringMd52(str, pass string) string {
	text := fmt.Sprintf("%s%s", StringMd5(str), pass)
	return StringMd5(text)
}

func IpToInt(ip net.IP) *big.Int {
	if v := ip.To4(); v != nil {
		return big.NewInt(0).SetBytes(v)
	}
	return big.NewInt(0).SetBytes(ip.To16())
}

func IntToIP(i *big.Int) net.IP {
	return net.IP(i.Bytes())
}

func StringToIP(i string) net.IP {
	return net.ParseIP(i).To4()
}

func GetIpAndPort(ip string) (string, string) {
	if strings.Contains(ip, ":") {
		arr := strings.Split(ip, ":")
		return arr[0], arr[1]
	}
	return ip, "22"
}

func Base64Encode(data string, a ...any) string {
	if len(a) > 0 {
		data = fmt.Sprintf(data, a...)
	}
	sEnc := base64.StdEncoding.EncodeToString([]byte(data))
	return fmt.Sprintf(sEnc)
}

func Base64Decode(data string) string {
	uDec, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return ""
	}
	return string(uDec)
}

// StringsContains 数组是否包含
func StringsContains(array []string, val string) (index int) {
	index = -1
	for i := 0; i < len(array); i++ {
		if array[i] == val {
			index = i
			return
		}
	}
	return
}

// InArray 元素是否存在数组中
func InArray(item string, items []string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

// Cmd 执行命令
func Cmd(arg ...string) (string, error) {
	output, err := exec.Command("/bin/bash", arg...).CombinedOutput()
	return string(output), err
}

// CmdSh 执行命令
func CmdSh(arg ...string) (string, error) {
	output, err := exec.Command("/bin/sh", arg...).CombinedOutput()
	return string(output), err
}

// CmdFile 执行命令
func CmdFile(filePath string) {
	_, _ = Cmd("-c", fmt.Sprintf("chmod +x %s", filePath))
	cmdString := exec.Command("/bin/bash", filePath)
	PrintCmdOutput(cmdString)
}

func PrintCmdOutput(cmd *exec.Cmd) {
	cmd.Stdin = os.Stdin

	var wg sync.WaitGroup
	wg.Add(2)
	//捕获标准输出
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("INFO:", err)
		os.Exit(1)
	}
	readout := bufio.NewReader(stdout)
	go func() {
		defer wg.Done()
		GetOutput(readout)
	}()

	//捕获标准错误
	stderr, err := cmd.StderrPipe()
	if err != nil {
		fmt.Println("ERROR:", err)
		os.Exit(1)
	}
	readerr := bufio.NewReader(stderr)
	go func() {
		defer wg.Done()
		GetOutput(readerr)
	}()

	//执行命令
	err = cmd.Run()
	if err != nil {
		return
	}
	wg.Wait()
}

func GetOutput(reader *bufio.Reader) {
	var sumOutput string //统计屏幕的全部输出内容
	outputBytes := make([]byte, 200)
	for {
		n, err := reader.Read(outputBytes) //获取屏幕的实时输出(并不是按照回车分割，所以要结合sumOutput)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println(err)
			sumOutput += err.Error()
		}
		output := string(outputBytes[:n])
		fmt.Print(output) //输出屏幕内容
		sumOutput += output
	}
}

// StructToJson 结构体转json
func StructToJson(data interface{}) string {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

// FormatEmail 格式化邮箱地址
func FormatEmail(email string) string {
	o, err := mail.ParseAddress(email)
	if err != nil {
		return ""
	}
	return o.Address
}

// IsEmail 判断是否是邮箱
func IsEmail(email string) bool {
	email = FormatEmail(email)
	if len(email) == 0 {
		return false
	}
	return true
}

// CheckOs 判断系统类型
func CheckOs() bool {
	return runtime.GOOS == "darwin" || runtime.GOOS == "linux"
}

// Test 正则判断
func Test(str, pattern string) bool {
	re := regexp.MustCompile(pattern)
	if re.MatchString(str) {
		return true
	} else {
		return false
	}
}

// RunDir 前面加上绝对路径
func RunDir(path string, a ...any) string {
	wd, _ := os.Getwd()
	if len(a) > 0 {
		path = fmt.Sprintf(path, a...)
	}
	return fmt.Sprintf("%s%s", wd, path)
}

// CacheDir 缓存目录
func CacheDir(path string, a ...any) string {
	return RunDir(fmt.Sprintf("/.cache%s", path), a...)
}
