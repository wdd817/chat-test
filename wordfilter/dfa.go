package wordfilter

// https://en.wikipedia.org/wiki/Trie
// Trie前缀树, 基于DFA进行的敏感词算法

type trieNode struct {
	isEndOfWord bool
	children map[byte]*trieNode
}

func newTrieNode() *trieNode {
	return &trieNode{
		isEndOfWord: false,
		children:    make(map[byte]*trieNode, 26),
	}
}

// Match index object
type matchIndex struct {
	start int
	end   int
}

func newMatchIndex(start, end int) *matchIndex {
	return &matchIndex{
		start: start,
		end:   end,
	}
}

// DFAUtil 敏感词查询过滤
type DFAUtil struct {
	root *trieNode
}

func (dfa *DFAUtil) insertWord(word []byte) {
	currNode := dfa.root
	for _, c := range word {
		if cildNode, exist := currNode.children[c]; !exist {
			cildNode = newTrieNode()
			currNode.children[c] = cildNode
			currNode = cildNode
		} else {
			currNode = cildNode
		}
	}

	currNode.isEndOfWord = true
}

func (dfa *DFAUtil) startsWith(prefix []byte) bool {
	currNode := dfa.root
	for _, c := range prefix {
		if childNode, exist := currNode.children[c]; !exist {
			return false
		} else {
			currNode = childNode
		}
	}

	return true
}

// 查询单个单词是否在trie树中
func (dfa *DFAUtil) searchWord(word []byte) bool {
	currNode := dfa.root
	for _, c := range word {
		if childNode, exist := currNode.children[c]; !exist {
			return false
		} else {
			currNode = childNode
		}
	}

	return currNode.isEndOfWord
}

// 查询一条文本, 返回多个匹配索引
func (dfa *DFAUtil) searchSentence(sentence string) (matchIndexList []*matchIndex) {
	start, end := 0, 1
	sentenceRuneList := []byte(sentence)

	startsWith := false
	for end <= len(sentenceRuneList) {
		// 逐个字母查询直至不属于trie树, 再判断子串中不存在敏感词
		if dfa.startsWith(sentenceRuneList[start:end]) {
			startsWith = true
			end += 1
		} else {
			if startsWith == true {
				for index := end - 1; index > start; index-- {
					if dfa.searchWord(sentenceRuneList[start:index]) {
						matchIndexList = append(matchIndexList, newMatchIndex(start, index-1))
						break
					}
				}
			}
			start, end = end-1, end+1
			startsWith = false
		}
	}

	// 如果整条文本都在trie中, 则再遍历一次所有子串
	if startsWith {
		for index := end - 1; index > start; index-- {
			if dfa.searchWord(sentenceRuneList[start:index]) {
				matchIndexList = append(matchIndexList, newMatchIndex(start, index-1))
				break
			}
		}
	}

	return
}

// 检查文本是否包含敏感词
func (dfa *DFAUtil) IsMatch(sentence string) bool {
	return len(dfa.searchSentence(sentence)) > 0
}

// 将传入文本中的敏感词用指定符号替换掉
func (dfa *DFAUtil) HandleWord(sentence string, replaceCh byte) (string, bool) {
	matchIndexList := dfa.searchSentence(sentence)
	if len(matchIndexList) == 0 {
		return sentence, false
	}

	// Manipulate
	sentenceList := []byte(sentence)
	for _, matchIndexObj := range matchIndexList {
		for index := matchIndexObj.start; index <= matchIndexObj.end; index++ {
			sentenceList[index] = replaceCh
		}
	}

	return string(sentenceList), true
}

// NewDFAUtil 创建DFA敏感词过滤器, 传入敏感词集合
func NewDFAUtil(wordList []string) *DFAUtil {
	this := &DFAUtil{
		root: newTrieNode(),
	}

	for _, word := range wordList {
		wordRuneList := []byte(word)
		if len(wordRuneList) > 0 {
			this.insertWord(wordRuneList)
		}
	}

	return this
}
