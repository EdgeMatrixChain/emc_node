package proof

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"github.com/emc-protocol/edge-matrix/helper/hex"
	"os/exec"
	"regexp"
	"testing"
	"time"
)

func Test_Unmarshal_PocRequestSuccessJson(t *testing.T) {
	resp := "{\"validator\":\"16Uiu2HAmKt7agigzA6oGDdMre4eCU7QER91vrW9M3aneiHEvGu1Y\",\"seed\":\"0x0e78238fd6e6686fd90f09df8c11c233763b2c4d79949818ee9f337001acc05c\"}"
	var obj struct {
		Validator string `json:"validator"`
		Seed      string `json:"seed"`
		Err       string `json:"err"`
	}

	if err := json.Unmarshal([]byte(resp), &obj); err != nil {
		t.Error(err)
	}
	if obj.Err != "" {
		t.Error(obj.Err)
	}
	t.Log("Validator:", obj.Validator)
	t.Log("Seed:", obj.Seed)
}

func Test_Unmarshal_PocRequestFailJson(t *testing.T) {

	resp := "{\"err\":\"block num too low\"}"
	var obj struct {
		Validator string `json:"validator"`
		Seed      string `json:"seed"`
		Err       string `json:"err"`
	}

	if err := json.Unmarshal([]byte(resp), &obj); err != nil {
		t.Error(err)
	}
	if obj.Err != "" {
		t.Log(obj.Err)
	}

}

func Test_HashJson2Map(t *testing.T) {
	testString := "[{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,8\",\"v\":\"5bb8f8e6c1758f54ce4ca4d8dbc703de2c670409dadf802c6a7f4c766c20a9d3\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,17\",\"v\":\"52ec3a0000ac9557eeedbb4516dfe49e237790d812939c438c8946adc45cc008\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,14\",\"v\":\"c566ee12a335ce6d4345d9f93e3ac4c10dfcb1c055460ecaef03459a034c9b59\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,18\",\"v\":\"b69a72f537fd9e7e98fccf19da805f4929b09ab5a44708f4da913438ea34b427\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,20\",\"v\":\"7e71c82280414697b16bc51493732175a255cd38a3675c463e1fe27e8b000587\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,23\",\"v\":\"1be7c86a028d13816479ac05de8ebb7429bdea0c0597c40899ee875a8291031f\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,0\",\"v\":\"6adced6ec51fc368614cd504ef4dc4ec2403e1131b61d336435aa657f0a8668c\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,2\",\"v\":\"2c1b3296d465cdb23e89f011d53be2b37e532072b34bdaeb2f529364a3d43885\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,4\",\"v\":\"35056fd4062fb0fec8f29a54e46e16412ffc39a67b641d5a52e71c60128de170\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,6\",\"v\":\"b82dbbc47301e956130ce32f02ecc7e8ed78610bd6c4c5d4c043b7b4ada6a9a2\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,36\",\"v\":\"aade79336bae3a525afc3444b0a6afb9dd6563bb1b199ec68e81851782cb3a04\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,57\",\"v\":\"423fe4aa95071ffc933949da6ec8d27cb5bd45efa0dbaac858aacde4664c9ef2\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,30\",\"v\":\"e7095f9602068f5f4743f6ce192ce794c65fac9e5a60429618300dd484c0326e\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,40\",\"v\":\"9b752e7e07c62ec647877ebcc8cd87a65e9a2e2216efa5db34293421e37f55b0\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,34\",\"v\":\"64249d9cbccb7f1c6b14f25317e9a7f4a94d7c0972956d00c273678169382fbb\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,46\",\"v\":\"8f22d75c931b6815f78bdbbfb8742c681b7ea20be45d6d33960df482aa69bc87\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,53\",\"v\":\"7c1643026a5b200a6b89f5c3e506424893b944b1359a60d8c5c9e96fab6dae02\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,54\",\"v\":\"cfdec623102ebcdcd147aad4cdea924f970bbb72487089faba5886e63df03882\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,1\",\"v\":\"c3cd9a929ce1d4a9dc7354fccdb52ff7eb80dc3ed0fc49151f4844d5ca6b0e0a\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,5\",\"v\":\"d7292aa823c4c0c5624514a085da1b387d7955907f11bc16c4372c6a87a2ce27\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,25\",\"v\":\"6ea55b43313f2d97a82f66530d55fbb2f9df37e163a6f309bdc93590a6d1e583\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,48\",\"v\":\"30734f406f86d1ccb398b3d7d7ae68f7ef2f00e7263b3bb44e09dcdcd479f428\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,19\",\"v\":\"0e8fa08e402c2835ccd50102c15138831e59df21e53b4696a28908d3057b3649\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,26\",\"v\":\"11d19188b4a76ed09c6874bf7553b80b1f5fbf7026b8c23b85795126dd76b1d4\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,37\",\"v\":\"82b75f273f50be26ca05100fa5b8f71b2271e82b88e9f285ee6a31e660558ff7\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,39\",\"v\":\"e5136a1bede4122199303999836b45326d8e548b4302fafb55a7dfc099b724d9\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,12\",\"v\":\"d9068216aed7f6d97e487d57a3cf0ba80981c73ee069789f768658b96b98085b\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,7\",\"v\":\"e03d0026ceefbe5be66a97bc2d68994aa0aad0654c56093c4a0b5fe5a3b2b31d\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,29\",\"v\":\"c6e01840d765edea0d21d1b6015aa7804dc323a055713b98f1633669c78e6d5f\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,45\",\"v\":\"d19f10d5da4f2c69db6dc73567487d4220eb8a0f0306f734e3eca0c970313324\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,49\",\"v\":\"7b91eeafc6ea4c16858f45374c793d7fbe505aee6167f4cef719a3fae44d7e0e\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,58\",\"v\":\"dfce75b8c8948ed67681e37d2b48c4708c9eef464c736082e424182e266f6309\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,35\",\"v\":\"43fb9efaec990b88b05449e165e611c3e3801b49bb223ceed760e08a656bc8ee\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,47\",\"v\":\"59004f5c22a06fd09aca8a8e68ab6a6cab3d0c89e11762dc46799d9099ce658e\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,50\",\"v\":\"cc1f4d2672b9b3d273aa6f5a8a13a28c8d2e9a00ff71ee41c8ca32cd8621f34e\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,51\",\"v\":\"9d63129e8da3ee42932ccb5f2cca74be81a0bca13889af89267c856f01a39dde\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,59\",\"v\":\"5da94e5f18887c68a92bee5d8db84cd28e5f1d48129b4c3db46eb08139be70a6\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,10\",\"v\":\"847aad0842a0f3149e8dd18a8e608488653f424604a3b3f64f4b33b31f5191a2\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,21\",\"v\":\"8e476b748c791fb1103157582bb1e097d9861ec1d10e95641104439efeef1f79\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,43\",\"v\":\"ce42edfe66adc65dd3b282dcd1053d947ffac31387f860e4f13fd489c96a0a8e\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,55\",\"v\":\"100032166842aefa83d35005347e8f7fd8e4392b1b3d5298f7c8fad68be689bb\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,11\",\"v\":\"28db95c9d3d94fc5141b372606180d5624781d682a8a04f1d954bde70977a187\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,16\",\"v\":\"5c413fec68473b2cc236f9590cec7661fc790b1a20435c970db6ce89a722c3c8\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,28\",\"v\":\"01844638964231beee02303eb70ab5bf0d6d8fe72710b5bdeff8e02811b02c7b\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,9\",\"v\":\"91b44b1ee782418f97492fd1f469319ff1364aa76f900c85617a6e667baddb3c\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,13\",\"v\":\"703c41cbb3b7d7005405be53928dc4650cf82c4dfa4f688e6cad2ed51be6700b\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,24\",\"v\":\"16887f3bf97eafbc3bf2e1aac515a86b099f0b008bb46957d5e05f8501053f7f\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,38\",\"v\":\"66012b6051a3230406c0ad12a970bb38c872313351d6e5cc2ac5c194b28adc69\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,44\",\"v\":\"70dbb96f2b28c6b11a24ea4ba5bc4a55b844cd3063e4c791347bd66f73d51e33\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,56\",\"v\":\"de6917d4bb2580b843ba2f7bfddd19666772bc9c4097acf1a97397e2d26c869b\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,15\",\"v\":\"cf3eca190d82321eea33f7b80308c0a935aea87546978428b63fdd81f0b1c558\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,32\",\"v\":\"8d325c91b1e777110198caf9ba47af7f30e9f8af451fd7539afdfc8e524e2287\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,33\",\"v\":\"3bfefa19eed80bc38633d1b3ffafb71ce2eadb94a2913eded5c23caff7579a08\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,42\",\"v\":\"105d0ea3571f0240939111b93765243860772923fe888800ed9d32382d0714ba\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,31\",\"v\":\"c1b5ccb03bae7f57f579277b414806a4b2c21147d6841d4129d35c6636efc86c\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,52\",\"v\":\"5c0dfa2a53d1af75857481e0d24fec465d3f4d302d737c21e930e18385ac6cc7\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,3\",\"v\":\"f7bc4f96df9f71ab67df5c10e0a1408786c0879a6c495942a43c1e9110413995\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,22\",\"v\":\"b9a4f8c0755310524bb1199008558c578cc0c01a8b6bc0d7401813c276b14658\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,27\",\"v\":\"5144d54634ae0e0e4a95d0c901e368f26a48e439e3dc7f739f386bf626f3d084\"}{\"k\":\"0xda37c9b99bd0d912acd7744d282a0bcdc3966cfb1934c37075cd025c05bf9050,41\",\"v\":\"04b49fe51245d20f160cded037788156b18c78f2e9a048ab3d93bca641z1afd87\"}]"
	fmt.Println(testString)

	var dataMapJson []map[string]string
	respBytes := []byte("[{\"k\":\"a\",\"v\":\"0e0e\"},{\"k\":\"b\",\"v\":\"0f0f\"}]")
	err := json.Unmarshal(respBytes, &dataMapJson)
	if err != nil {
		t.Error(err)
		return
	}

	dataMap := make(map[string][]byte)
	for _, data := range dataMapJson {
		bytes, err := hex.DecodeString(data["v"])
		if err != nil {
			t.Error(err)
			return
		}
		t.Log("k:", data["k"], ",v:", data["v"])
		dataMap[data["k"]] = bytes
	}
	t.Log(dataMap)
}
func TestProofBlockNumber(t *testing.T) {
	nums := make([]int64, 100)
	i := int64(0)
	for i < 100 {
		nums[i] = 1892345 + i
		i += 1
	}
	for _, blockNumber := range nums {
		t.Log("blockNumber:", blockNumber)
		blockNumberFixed := (blockNumber / 30) * 30
		t.Log("blockNumberFixed:", blockNumberFixed)
	}

}

func TestProofByHash(t *testing.T) {
	var data = make(map[string]*[]byte)
	target := "0000"
	loops := 60
	i := 0
	start := time.Now()
	for i < loops {
		randBytes := make([]byte, 32)
		_, err := rand.Read(randBytes)
		if err != nil {
			return
		}
		seed := hex.EncodeToHex(randBytes)
		_, bytes, err := ProofByCalcHash(seed, target, time.Second*3)
		if err != nil {
			t.Log(fmt.Sprintf("err: %s", err.Error()))
			return
		}
		data[seed] = &bytes
		i += 1
	}
	t.Log(fmt.Sprintf("calc time			: %fs", time.Since(start).Seconds()))

	validateSuccess := 0
	validateStart := time.Now()
	for seed, bytes := range data {
		validateHash := ValidateHash(seed, target, *bytes)
		if validateHash {
			validateSuccess += 1
		}
	}
	t.Log(fmt.Sprintf("validate time		: %dms", time.Since(validateStart).Microseconds()))
	t.Log(fmt.Sprintf("validate success	: %d/%d", validateSuccess, loops))

}

func TestCpuInfo(t *testing.T) {
	cmd := exec.Command("wmic", "cpu", "get", "ProcessorID")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(out))
	str := string(out)
	reg := regexp.MustCompile("\\s+")
	str = reg.ReplaceAllString(str, "")
	t.Log(str)
}
