package dgraph_helper

//
//func TestMutateAssembler_Common(t *testing.T) {
//	stream := newFakeMutationStream()
//	l1 := NewLedger("ledgerhash1")
//	l2 := NewLedger("ledgerhash2")
//	stream.AddData(
//		l1, l2,
//		NewBlock("blockhash0", 0, l1),
//		NewBlock("blockahsh1", 1, l1),
//		NewBlock("blockhash2", 1, l2),
//		NewBlock("blockhash3", 2, l2),
//	)
//	fakeDB := newFakeDb()
//	cache := NewUidLruCache(fakeDB, 100)
//	assembler := NewMutateAssembler(stream, cache)
//	result, err := assembler.Output()
//	assert.Nil(t, err)
//	t.Log(result)
//}
//
//func TestMutateAssembler_Link(t *testing.T) {
//	stream := newFakeMutationStream()
//	l1 := newLedgerBlockLink("ledgerhash1", "blockhash1")
//	l2 := newLedgerBlockLink("ledgerhash2", "blockhash2")
//	stream.AddData(
//		l1, l2,
//	)
//	fakeDB := newFakeDb()
//	fakeDB.SetUID("block-hash_id", "blockhash1", "0x002")
//	fakeDB.SetUID("ledger-hash_id", "ledgerhash2", "0x003")
//	fakeDB.SetUID("block-hash_id", "blockhash2", "0x004")
//	cache := NewUidLruCache(fakeDB, 100)
//	cache.UpdateUid("ledger-hash_id||ledgerhash1", "0x001")
//	assembler := NewMutateAssembler(stream, cache)
//	result, err := assembler.Output()
//	assert.Nil(t, err)
//	t.Log(result)
//}
//
//func newFakeDb() *FakeDb {
//	return &FakeDb{
//		kvs: map[string]string{},
//	}
//}
//
//type FakeDb struct {
//	kvs map[string]string
//}
//
//func (db *FakeDb) SetUID(predict, value, uid string) {
//	db.kvs[predict+"-"+value] = uid
//}
//
//func (db *FakeDb) QueryUID(predict, value string) (uid string, exists bool, e error) {
//	uid, exists = db.kvs[predict+"-"+value]
//	fmt.Printf("get uid from db by %s-%s -> %s\n", predict, value, uid)
//	return
//}
//
//func newFakeMutationStream() *FakeMutationStream {
//	return &FakeMutationStream{
//		list: list.New(),
//	}
//}
//
//type FakeMutationStream struct {
//	list *list.List
//}
//
//func (stream *FakeMutationStream) AddData(mds ...MutationData) {
//	for _, md := range mds {
//		stream.list.PushBack(md)
//	}
//}
//
//func (stream *FakeMutationStream) Next() (data MutationData, eof bool, e error) {
//	if stream.list.Len() <= 0 {
//		eof = true
//		return
//	}
//	front := stream.list.Front()
//	data = front.Value.(MutationData)
//	stream.list.Remove(front)
//	return
//}
//
//func newLedgerBlockLink(ledger, block string) *LedgerBlockLink {
//	return &LedgerBlockLink{
//		ledger: ledger,
//		block:  block,
//	}
//}
//
//type LedgerBlockLink struct {
//	ledger    string
//	ledgerUID string
//	block     string
//	blockUID  string
//}
//
//func (link *LedgerBlockLink) SetUidLeft(uid string) {
//	link.ledgerUID = uid
//}
//
//func (link *LedgerBlockLink) SetUidRight(uid string) {
//	link.blockUID = uid
//}
//
//func (link *LedgerBlockLink) LeftQueryBy() (predict, value string, ok bool) {
//	predict = "ledger-hash_id"
//	value = link.ledger
//	ok = true
//	return
//}
//
//func (link *LedgerBlockLink) RightQueryBy() (predict, value string, ok bool) {
//	predict = "block-hash_id"
//	value = link.block
//	ok = true
//	return
//}
//
//func (link *LedgerBlockLink) Mutations() (mutations Mutations) {
//	mutations = mutations.Add(
//		NewMutation(
//			MutationItemUid(link.ledgerUID),
//			MutationItemUid(link.blockUID),
//			MutationPredict("ledger-block"),
//		),
//	)
//	return
//}
//
//type Ledger struct {
//	Uid  string
//	Hash string
//}
//
//func NewLedger(hashID string) *Ledger {
//	return &Ledger{
//		Hash: hashID,
//	}
//}
//
//func (ledger *Ledger) MutationName() string {
//	return "ledger-" + ledger.Hash
//}
//
//func (ledger *Ledger) Mutations() (mutations Mutations) {
//	mutations = mutations.Add(
//		NewMutation(
//			MutationItemEmpty(ledger.MutationName()),
//			MutationItemValue(ledger.Hash),
//			MutationPredict("ledger-hash_id"),
//		),
//	)
//	return
//}
//
//type Block struct {
//	Uid    string
//	Hash   string
//	Height int64
//	Ledger *Ledger
//}
//
//func NewBlock(hashID string, height int64, ledger *Ledger) *Block {
//	return &Block{
//		Hash:   hashID,
//		Height: height,
//		Ledger: ledger,
//	}
//}
//
//func (block *Block) UniqueMutationName() string {
//	return fmt.Sprintf("block%s", block.Hash)
//}
//
//func (block *Block) Mutations() (mutations Mutations) {
//	mutations = mutations.Add(
//		NewMutation(
//			MutationItemEmpty(block.UniqueMutationName()),
//			MutationItemValue(block.Hash),
//			MutationPredict("block-hash_id"),
//		),
//		NewMutation(
//			MutationItemEmpty(block.UniqueMutationName()),
//			MutationItemValue(strconv.FormatInt(block.Height, 10)),
//			MutationPredict("block-height"),
//		),
//	)
//	return
//}
