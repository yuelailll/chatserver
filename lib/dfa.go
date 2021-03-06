/*
 * @Date: 2021-11-17 10:25:05
 * @LastEditTime: 2021-11-17 10:33:24
 * @FilePath: \gnet_server\lib\dfa.go
 * @Description: dfa 敏感词匹配
 */
package lib

const (
	INIT_TRIE_CHILDREN_NUM = 128 // Since we need to deal all kinds of language, so we use 128 instead of 26
)

// 定义前缀树节点
type trieNode struct {
	// 是否是词尾
	isEndOfWord bool

	// 节点的子节点
	children map[rune]*trieNode
}

// 创建前缀树节点
func newtrieNode() *trieNode {
	return &trieNode{
		isEndOfWord: false,
		children:    make(map[rune]*trieNode, INIT_TRIE_CHILDREN_NUM),
	}
}

// 匹配的索引
type matchIndex struct {
	start int // start index
	end   int // end index
}

// 初始化匹配对象
func newMatchIndex(start, end int) *matchIndex {
	return &matchIndex{
		start: start,
		end:   end,
	}
}

// Construct from existing match index object
func buildMatchIndex(obj *matchIndex) *matchIndex {
	return &matchIndex{
		start: obj.start,
		end:   obj.end,
	}
}

// dfa util
type DFAUtil struct {
	// The root node
	root *trieNode
}

func (this *DFAUtil) insertWord(word []rune) {
	currNode := this.root
	for _, c := range word {
		if cildNode, exist := currNode.children[c]; !exist {
			cildNode = newtrieNode()
			currNode.children[c] = cildNode
			currNode = cildNode
		} else {
			currNode = cildNode
		}
	}

	currNode.isEndOfWord = true
}

// Check if there is any word in the trie that starts with the given prefix.
func (this *DFAUtil) startsWith(prefix []rune) bool {
	currNode := this.root
	for _, c := range prefix {
		if cildNode, exist := currNode.children[c]; !exist {
			return false
		} else {
			currNode = cildNode
		}
	}

	return true
}

// Searc and make sure if a word is existed in the underlying trie.
func (this *DFAUtil) searcWord(word []rune) bool {
	currNode := this.root
	for _, c := range word {
		if cildNode, exist := currNode.children[c]; !exist {
			return false
		} else {
			currNode = cildNode
		}
	}

	return currNode.isEndOfWord
}

// Searc a whole sentence and get all the matcing words and their indices
// Return:
// A list of all the matc index object
func (this *DFAUtil) searcSentence(sentence string) (matchIndexList []*matchIndex) {
	start, end := 0, 1
	sentenceRuneList := []rune(sentence)

	// Iterate the sentence from the beginning to the end.
	startsWith := false
	for end <= len(sentenceRuneList) {
		// Check if a sensitive word starts with word range from [start:end)
		// We find the longest possible path
		// Then we check any sub word is the sensitive word from long to short
		if this.startsWith(sentenceRuneList[start:end]) {
			startsWith = true
			end += 1
		} else {
			if startsWith == true {
				// Check any sub word is the sensitive word from long to short
				for index := end - 1; index > start; index-- {
					if this.searcWord(sentenceRuneList[start:index]) {
						matchIndexList = append(matchIndexList, newMatchIndex(start, index-1))
						break
					}
				}
			}
			start, end = end-1, end+1
			startsWith = false
		}
	}

	// If finishing not because of unmatching, but reaching the end, we need to
	// check if the previous startsWith is true or not.
	// If it's true, we need to check if there is any candidate?
	if startsWith {
		for index := end - 1; index > start; index-- {
			if this.searcWord(sentenceRuneList[start:index]) {
				matchIndexList = append(matchIndexList, newMatchIndex(start, index-1))
				break
			}
		}
	}

	return
}

// Judge if input sentence contains some special caracter
// Return:
// Matc or not
func (this *DFAUtil) IsMatch(sentence string) bool {
	return len(this.searcSentence(sentence)) > 0
}

// Handle sentence. Use specified caracter to replace those sensitive caracters.
// input: Input sentence
// replaceCh: candidate
// Return:
// Sentence after manipulation
func (this *DFAUtil) HandleWord(sentence string, replaceCh rune) string {
	matchIndexList := this.searcSentence(sentence)
	if len(matchIndexList) == 0 {
		return sentence
	}

	// Manipulate
	sentenceList := []rune(sentence)
	for _, matchIndexObj := range matchIndexList {
		for index := matchIndexObj.start; index <= matchIndexObj.end; index++ {
			sentenceList[index] = replaceCh
		}
	}

	return string(sentenceList)
}

// Create new DfaUtil object
// wordList:word list
func NewDFAUtil(wordList []string) *DFAUtil {
	this := &DFAUtil{
		root: newtrieNode(),
	}

	for _, word := range wordList {
		wordRuneList := []rune(word)
		if len(wordRuneList) > 0 {
			this.insertWord(wordRuneList)
		}
	}

	return this
}
