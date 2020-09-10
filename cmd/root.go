package cmd

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

type Kind uint

const (
	KindDigit = iota + 1
	KindAlpha
	KindPunct
)

type ShiftState uint

const (
	ShiftStateOn = iota + 1
	ShiftStateOff
)

type KeyLine struct {
	UpperChars string
	LowerChars string
}

func (kl *KeyLine) size() int {
	return len(kl.LowerChars)
}

func (kl *KeyLine) At(index int) []byte {
	res := make([]byte, 0, 2)
	if index < 0 || index >= kl.size() {
		return res
	}
	if kl.LowerChars[index] != ' ' {
		res = append(res, kl.LowerChars[index])
	}
	if kl.UpperChars[index] != ' ' {
		res = append(res, kl.UpperChars[index])
	}
	return res
}

func (kl *KeyLine) Around(index int) []byte {
	res := make([]byte, 0, 4)
	toCheck := []string{kl.UpperChars, kl.LowerChars}
	for _, str := range toCheck {
		if index == 0 && str[index+1] != ' ' {
			res = append(res, str[index+1])
		} else if index == len(str)-1 && str[index-1] != ' ' {
			res = append(res, str[index-1])
		} else {
			if str[index-1] != ' ' {
				res = append(res, str[index-1])
			}
			if str[index+1] != ' ' {
				res = append(res, str[index+1])
			}
		}
	}
	return res
}

type KeyMatrix struct {
	lines       []*KeyLine
	illegalNext map[byte]map[byte]bool
}

func NewQwertyMatrix() *KeyMatrix {
	m := &KeyMatrix{
		lines: []*KeyLine{
			&KeyLine{
				UpperChars: "~!@#$%^&*()_+",
				LowerChars: "`1234567890-=",
			},
			&KeyLine{
				UpperChars: ` QWERTYUIOP{}|`,
				LowerChars: ` qwertyuiop[]\`,
			},
			&KeyLine{
				UpperChars: ` ASDFGHJKL:"`,
				LowerChars: ` asdfghjkl;'`,
			},
			&KeyLine{
				UpperChars: ` ZXCVBNM<>?`,
				LowerChars: ` zxcvbnm,./`,
			},
		},
		illegalNext: make(map[byte]map[byte]bool),
	}
	m.init()
	return m
}

func (km *KeyMatrix) init() {
	for i := 0; i < km.size(); i += 1 {
		for j := 0; j < km.lines[i].size(); j += 1 {
			around := km.Around(i, j)
			chars := km.lines[i].At(j)
			for _, c := range chars {
				s, ok := km.illegalNext[c]
				if !ok {
					s = make(map[byte]bool)
					km.illegalNext[c] = s
				}
				for _, ch := range around {
					s[ch] = true
				}
				// treat all upper/lower chars on curruent key as neighbour
				for _, ch := range chars {
					s[ch] = true
				}
			}
		}
	}
}

func (km *KeyMatrix) IsIllegialNext(ch1 byte, ch2 byte) bool {
	if next, ok := km.illegalNext[ch1]; ok {
		if _, ok := next[ch2]; ok {
			return true
		}
	}
	return false
}

func (km *KeyMatrix) size() int {
	return len(km.lines)
}

func (km *KeyMatrix) Around(row int, col int) []byte {
	if row == 0 {
		res := append(km.lines[row].Around(col), km.lines[row+1].At(col)...)
		return res
	} else if row == km.size()-1 {
		res := append(km.lines[row].Around(col), km.lines[row-1].At(col)...)
		return res
	} else {
		res := append(km.lines[row].Around(col), km.lines[row-1].At(col)...)
		res = append(res, km.lines[row+1].At(col)...)
		return res
	}
}

const (
	strDigit = "0123456789"
	strAlpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	strPunct = `~!@#$%^&*()_+[]{}|\;:'",./<>?`
)

var (
	printVersion bool
	useDigit     bool
	useAlpha     bool
	usePunct     bool
	strictMode   bool
	length       uint
	digitNum     uint
	punctNum     uint
	qwertyMatrix *KeyMatrix
)

func run() {
	if printVersion {
		printVersionInfo()
		os.Exit(0)
	}

	if length == 0 {
		os.Exit(0)
	}

	if strictMode {
		useDigit = true
		useAlpha = true
		usePunct = true
	}
	kinds := make([]Kind, 0, 3)
	alphaNum := length
	if useDigit && digitNum > 0 {
		if length < digitNum {
			fmt.Printf("invalid length: %d < digitNum: %d", length, digitNum)
			os.Exit(1)
		}
		kinds = append(kinds, KindDigit)
		alphaNum -= digitNum
	}
	if usePunct && punctNum > 0 {
		if length < punctNum {
			fmt.Printf("invalid length: %d < punctNum: %d", length, punctNum)
			os.Exit(1)
		}
		if alphaNum < punctNum {
			fmt.Printf("invalid length: %d < digitNum: %d + punctNum: %d", length, digitNum, punctNum)
			os.Exit(1)
		}
		kinds = append(kinds, KindPunct)
		alphaNum -= punctNum
	}
	if alphaNum > 0 {
		kinds = append(kinds, KindAlpha)
	}

	var res bytes.Buffer
	var last, next byte
	var nextKind int
	var dn, pn, an uint
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := uint(0); i < length; i += 1 {
	INNER:
		for {
			nextKind = r.Intn(len(kinds))
			switch kinds[nextKind] {
			case KindDigit:
				next = strDigit[r.Intn(len(strDigit))]
			case KindPunct:
				next = strPunct[r.Intn(len(strPunct))]
			case KindAlpha:
				next = strAlpha[r.Intn(len(strAlpha))]
			default:
				panic("should not go here")
			}
			if !qwertyMatrix.IsIllegialNext(last, next) {
				break INNER
			}
		}
		res.WriteByte(next)
		last = next
		switch kinds[nextKind] {
		case KindDigit:
			dn += 1
		case KindPunct:
			pn += 1
		case KindAlpha:
			an += 1
		default:
			panic("should not go here")
		}
		kinds = updateKinds(alphaNum, kinds, dn, pn, an)
	}
	fmt.Printf("%s\n", res.String())
}

func updateKinds(alphaNum uint, kinds []Kind, dn, pn, an uint) []Kind {
	if useDigit && dn >= digitNum {
		kinds = filterOut(kinds, KindDigit)
	}
	if usePunct && pn >= punctNum {
		kinds = filterOut(kinds, KindPunct)
	}
	if an >= alphaNum {
		kinds = filterOut(kinds, KindAlpha)
	}
	return kinds
}

func filterOut(kinds []Kind, out Kind) []Kind {
	res := make([]Kind, 0, len(kinds))
	for _, k := range kinds {
		if k == out {
			continue
		}
		res = append(res, k)
	}
	return res
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mkpw",
	Short: "Make password",
	Long:  `Make password`,
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initialize)

	rootCmd.PersistentFlags().BoolVar(&printVersion, "version", false,
		"print version information and exit")
	rootCmd.PersistentFlags().BoolVar(&useDigit, "digit", true, "use digit characters")
	rootCmd.PersistentFlags().BoolVar(&useAlpha, "alpha", true, "use alphabet characters")
	rootCmd.PersistentFlags().BoolVar(&usePunct, "punct", false, "use punctuation characters")
	rootCmd.PersistentFlags().BoolVar(&strictMode, "strict", false, "switch to strict mode")
	rootCmd.PersistentFlags().UintVar(&length, "length", 8, "the total length")
	rootCmd.PersistentFlags().UintVar(&digitNum, "digitNum", 2, "the digit char number")
	rootCmd.PersistentFlags().UintVar(&punctNum, "punctNum", 1, "the punctuation char number")
}

func initialize() {
	qwertyMatrix = NewQwertyMatrix()
}
