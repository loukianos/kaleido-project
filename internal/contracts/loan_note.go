// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package contracts

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// LoanNoteLoanTerms is an auto generated low-level Go binding around an user-defined struct.
type LoanNoteLoanTerms struct {
	Principal *big.Int
	AprBps    uint16
	Maturity  uint64
	Status    uint8
}

// LoanNoteMetaData contains all meta data concerning the LoanNote contract.
var LoanNoteMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"admin\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"AccessControlBadConfirmation\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"neededRole\",\"type\":\"bytes32\"}],\"name\":\"AccessControlUnauthorizedAccount\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_fromTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_toTokenId\",\"type\":\"uint256\"}],\"name\":\"BatchMetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"aprBps\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maturity\",\"type\":\"uint64\"}],\"name\":\"LoanOriginated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"LoanSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumLoanNote.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"LoanStatusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"MetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"previousAdminRole\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"newAdminRole\",\"type\":\"bytes32\"}],\"name\":\"RoleAdminChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleGranted\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"RoleRevoked\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"DEFAULT_ADMIN_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ORIGINATOR_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"SERVICER_ROLE\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"}],\"name\":\"getRoleAdmin\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"grantRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"hasRole\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"markDefaulted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"lender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"aprBps\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"maturity\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"tokenURI_\",\"type\":\"string\"}],\"name\":\"originate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"callerConfirmation\",\"type\":\"address\"}],\"name\":\"renounceRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"role\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"revokeRole\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"terms\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"aprBps\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"maturity\",\"type\":\"uint64\"},{\"internalType\":\"enumLoanNote.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"internalType\":\"structLoanNote.LoanTerms\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b50604051611dc8380380611dc883398101604081905261002e916101a0565b604051806040016040528060088152602001674c6f616e4e6f746560c01b815250604051806040016040528060048152602001632627a0a760e11b815250815f908161007a9190610265565b5060016100878282610265565b5061009691505f9050826100f3565b506100c17f59abfac6520ec36a6556b2a4dd949cc40007459bcd5cd2507f1e5cc77b6bc97e826100f3565b506100ec7f250b76734a070a69c7b3930477dd35007ad9c9d0952e97903fdafb2db6980537826100f3565b505061031f565b5f8281526007602090815260408083206001600160a01b038516845290915281205460ff16610197575f8381526007602090815260408083206001600160a01b03861684529091529020805460ff1916600117905561014f3390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a450600161019a565b505f5b92915050565b5f602082840312156101b0575f5ffd5b81516001600160a01b03811681146101c6575f5ffd5b9392505050565b634e487b7160e01b5f52604160045260245ffd5b600181811c908216806101f557607f821691505b60208210810361021357634e487b7160e01b5f52602260045260245ffd5b50919050565b601f82111561026057805f5260205f20601f840160051c8101602085101561023e5750805b601f840160051c820191505b8181101561025d575f815560010161024a565b50505b505050565b81516001600160401b0381111561027e5761027e6101cd565b6102928161028c84546101e1565b84610219565b6020601f8211600181146102c4575f83156102ad5750848201515b5f19600385901b1c1916600184901b17845561025d565b5f84815260208120601f198516915b828110156102f357878501518255602094850194600190920191016102d3565b508482101561031057868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b611a9c8061032c5f395ff3fe608060405234801561000f575f5ffd5b506004361061016d575f3560e01c806373216450116100d9578063a22cb46511610093578063c87b56dd1161006e578063c87b56dd1461036c578063d547741f1461037f578063daf18fba14610392578063e985e9c5146103a5575f5ffd5b8063a22cb46514610326578063b88d4fde14610339578063c0aa0e8a1461034c575f5ffd5b806373216450146102b75780638df82800146102ca57806391d14854146102dd57806395d89b41146102f0578063a125142b146102f8578063a217fddf1461031f575f5ffd5b8063248a9ca31161012a578063248a9ca3146102365780632f2ff15d1461025857806336568abe1461026b57806342842e0e1461027e5780636352211e1461029157806370a08231146102a4575f5ffd5b806301ffc9a71461017157806306fdde0314610199578063081812fc146101ae578063095ea7b3146101d957806309a01608146101ee57806323b872dd14610223575b5f5ffd5b61018461017f3660046114f0565b6103b8565b60405190151581526020015b60405180910390f35b6101a16103c8565b6040516101909190611539565b6101c16101bc36600461154b565b610457565b6040516001600160a01b039091168152602001610190565b6101ec6101e736600461157d565b61047e565b005b6102157f250b76734a070a69c7b3930477dd35007ad9c9d0952e97903fdafb2db698053781565b604051908152602001610190565b6101ec6102313660046115a5565b61048d565b61021561024436600461154b565b5f9081526007602052604090206001015490565b6101ec6102663660046115df565b61051b565b6101ec6102793660046115df565b61053f565b6101ec61028c3660046115a5565b610577565b6101c161029f36600461154b565b610591565b6102156102b2366004611609565b61059b565b6101ec6102c536600461154b565b6105e0565b6101ec6102d836600461154b565b610674565b6101846102eb3660046115df565b61073c565b6101a1610766565b6102157f59abfac6520ec36a6556b2a4dd949cc40007459bcd5cd2507f1e5cc77b6bc97e81565b6102155f81565b6101ec610334366004611622565b610775565b6101ec61034736600461166f565b610780565b61035f61035a36600461154b565b610798565b604051610190919061177d565b6101a161037a36600461154b565b610848565b6101ec61038d3660046115df565b610953565b6102156103a03660046117c2565b610977565b6101846103b336600461187d565b610b20565b5f6103c282610b4d565b92915050565b60605f80546103d6906118a5565b80601f0160208091040260200160405190810160405280929190818152602001828054610402906118a5565b801561044d5780601f106104245761010080835404028352916020019161044d565b820191905f5260205f20905b81548152906001019060200180831161043057829003601f168201915b5050505050905090565b5f61046182610b71565b505f828152600460205260409020546001600160a01b03166103c2565b610489828233610ba9565b5050565b6001600160a01b0382166104bb57604051633250574960e11b81525f60048201526024015b60405180910390fd5b5f6104c7838333610bb6565b9050836001600160a01b0316816001600160a01b031614610515576040516364283d7b60e01b81526001600160a01b03808616600483015260248201849052821660448201526064016104b2565b50505050565b5f8281526007602052604090206001015461053581610ca8565b6105158383610cb5565b6001600160a01b03811633146105685760405163334bd91960e11b815260040160405180910390fd5b6105728282610d46565b505050565b61057283838360405180602001604052805f815250610780565b5f6103c282610b71565b5f6001600160a01b0382166105c5576040516322718ad960e21b81525f60048201526024016104b2565b506001600160a01b03165f9081526003602052604090205490565b7f250b76734a070a69c7b3930477dd35007ad9c9d0952e97903fdafb2db698053761060a81610ca8565b61061382610b71565b505f8281526009602052604090819020600101805460ff60501b1916600160511b1790555182907f9d0b8f6161220422fcfcf3cbe3b12d5148060bea52d7d74395488cae75d2e46f90610668906002906118dd565b60405180910390a25050565b7f250b76734a070a69c7b3930477dd35007ad9c9d0952e97903fdafb2db698053761069e81610ca8565b6106a782610b71565b50817f9d0b8f6161220422fcfcf3cbe3b12d5148060bea52d7d74395488cae75d2e46f60016040516106d991906118dd565b60405180910390a260405182907ff2257195ce09a74d98ff0578eba9d4573bc34bb4b21473e363535fe3485e0efa905f90a261071482610db1565b505f90815260096020526040812090815560010180546affffffffffffffffffffff19169055565b5f9182526007602090815260408084206001600160a01b0393909316845291905290205460ff1690565b6060600180546103d6906118a5565b610489338383610de9565b61078b84848461048d565b6105153385858585610e87565b6107bf604080516080810182525f8082526020820181905291810182905290606082015290565b6107c882610b71565b505f82815260096020908152604091829020825160808101845281548152600182015461ffff8116938201939093526201000083046001600160401b031693810193909352906060830190600160501b900460ff16600281111561082e5761082e611749565b600281111561083f5761083f611749565b90525092915050565b606061085382610b71565b505f828152600660205260408120805461086c906118a5565b80601f0160208091040260200160405190810160405280929190818152602001828054610898906118a5565b80156108e35780601f106108ba576101008083540402835291602001916108e3565b820191905f5260205f20905b8154815290600101906020018083116108c657829003601f168201915b505050505090505f6108ff60408051602081019091525f815290565b905080515f03610910575092915050565b81511561094257808260405160200161092a929190611902565b60405160208183030381529060405292505050919050565b61094b84610faf565b949350505050565b5f8281526007602052604090206001015461096d81610ca8565b6105158383610d46565b5f7f59abfac6520ec36a6556b2a4dd949cc40007459bcd5cd2507f1e5cc77b6bc97e6109a281610ca8565b60088054905f6109b183611916565b91905055915060405180608001604052808881526020018761ffff168152602001866001600160401b031681526020015f60028111156109f3576109f3611749565b90525f838152600960209081526040918290208351815590830151600182018054938501516001600160401b0316620100000269ffffffffffffffffffff1990941661ffff9092169190911792909217808355606084015191929060ff60501b1916600160501b836002811115610a6c57610a6c611749565b0217905550905050610a7e8883611020565b610abd8285858080601f0160208091040260200160405190810160405280939291908181526020018383808284375f9201919091525061103992505050565b6040805188815261ffff881660208201526001600160401b0387168183015290516001600160a01b038a169184917f3ade124d39b4fb6fb328df4ccd68907e63c76c84fe020c179e2195101b4644009181900360600190a3509695505050505050565b6001600160a01b039182165f90815260056020908152604080832093909416825291909152205460ff1690565b5f6001600160e01b03198216637965db0b60e01b14806103c257506103c282611088565b5f818152600260205260408120546001600160a01b0316806103c257604051637e27328960e01b8152600481018490526024016104b2565b61057283838360016110ac565b5f828152600260205260408120546001600160a01b0390811690831615610be257610be28184866111b0565b6001600160a01b03811615610c1c57610bfd5f855f5f6110ac565b6001600160a01b0381165f90815260036020526040902080545f190190555b6001600160a01b03851615610c4a576001600160a01b0385165f908152600360205260409020805460010190555b5f8481526002602052604080822080546001600160a01b0319166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b610cb28133611214565b50565b5f610cc0838361073c565b610d3f575f8381526007602090815260408083206001600160a01b03861684529091529020805460ff19166001179055610cf73390565b6001600160a01b0316826001600160a01b0316847f2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d60405160405180910390a45060016103c2565b505f6103c2565b5f610d51838361073c565b15610d3f575f8381526007602090815260408083206001600160a01b0386168085529252808320805460ff1916905551339286917ff6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b9190a45060016103c2565b5f610dbd5f835f610bb6565b90506001600160a01b03811661048957604051637e27328960e01b8152600481018390526024016104b2565b6001600160a01b038216610e1b57604051630b61174360e31b81526001600160a01b03831660048201526024016104b2565b6001600160a01b038381165f81815260056020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b6001600160a01b0383163b15610fa857604051630a85bd0160e11b81526001600160a01b0384169063150b7a0290610ec990889088908790879060040161193a565b6020604051808303815f875af1925050508015610f03575060408051601f3d908101601f19168201909252610f0091810190611976565b60015b610f6a573d808015610f30576040519150601f19603f3d011682016040523d82523d5f602084013e610f35565b606091505b5080515f03610f6257604051633250574960e11b81526001600160a01b03851660048201526024016104b2565b805181602001fd5b6001600160e01b03198116630a85bd0160e11b14610fa657604051633250574960e11b81526001600160a01b03851660048201526024016104b2565b505b5050505050565b6060610fba82610b71565b505f610fd060408051602081019091525f815290565b90505f815111610fee5760405180602001604052805f815250611019565b80610ff88461124d565b604051602001611009929190611902565b6040516020818303038152906040525b9392505050565b610489828260405180602001604052805f8152506112dc565b5f82815260066020526040902061105082826119d5565b506040518281527ff8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce79060200160405180910390a15050565b5f6001600160e01b03198216632483248360e11b14806103c257506103c2826112f3565b80806110c057506001600160a01b03821615155b15611181575f6110cf84610b71565b90506001600160a01b038316158015906110fb5750826001600160a01b0316816001600160a01b031614155b801561110e575061110c8184610b20565b155b156111375760405163a9fbf51f60e01b81526001600160a01b03841660048201526024016104b2565b811561117f5783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b50505f90815260046020526040902080546001600160a01b0319166001600160a01b0392909216919091179055565b6111bb838383611342565b610572576001600160a01b0383166111e957604051637e27328960e01b8152600481018290526024016104b2565b60405163177e802f60e01b81526001600160a01b0383166004820152602481018290526044016104b2565b61121e828261073c565b6104895760405163e2517d3f60e01b81526001600160a01b0382166004820152602481018390526044016104b2565b60605f611259836113a3565b60010190505f816001600160401b038111156112775761127761165b565b6040519080825280601f01601f1916602001820160405280156112a1576020820181803683370190505b5090508181016020015b5f19016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a85049450846112ab57509392505050565b6112e6838361147a565b610572335f858585610e87565b5f6001600160e01b031982166380ac58cd60e01b148061132357506001600160e01b03198216635b5e139f60e01b145b806103c257506301ffc9a760e01b6001600160e01b03198316146103c2565b5f6001600160a01b0383161580159061094b5750826001600160a01b0316846001600160a01b0316148061137b575061137b8484610b20565b8061094b5750505f908152600460205260409020546001600160a01b03908116911614919050565b5f8072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b83106113e15772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6d04ee2d6d415b85acef8100000000831061140d576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc10000831061142b57662386f26fc10000830492506010015b6305f5e1008310611443576305f5e100830492506008015b612710831061145757612710830492506004015b60648310611469576064830492506002015b600a83106103c25760010192915050565b6001600160a01b0382166114a357604051633250574960e11b81525f60048201526024016104b2565b5f6114af83835f610bb6565b90506001600160a01b03811615610572576040516339e3563760e11b81525f60048201526024016104b2565b6001600160e01b031981168114610cb2575f5ffd5b5f60208284031215611500575f5ffd5b8135611019816114db565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f611019602083018461150b565b5f6020828403121561155b575f5ffd5b5035919050565b80356001600160a01b0381168114611578575f5ffd5b919050565b5f5f6040838503121561158e575f5ffd5b61159783611562565b946020939093013593505050565b5f5f5f606084860312156115b7575f5ffd5b6115c084611562565b92506115ce60208501611562565b929592945050506040919091013590565b5f5f604083850312156115f0575f5ffd5b8235915061160060208401611562565b90509250929050565b5f60208284031215611619575f5ffd5b61101982611562565b5f5f60408385031215611633575f5ffd5b61163c83611562565b915060208301358015158114611650575f5ffd5b809150509250929050565b634e487b7160e01b5f52604160045260245ffd5b5f5f5f5f60808587031215611682575f5ffd5b61168b85611562565b935061169960208601611562565b92506040850135915060608501356001600160401b038111156116ba575f5ffd5b8501601f810187136116ca575f5ffd5b80356001600160401b038111156116e3576116e361165b565b604051601f8201601f19908116603f011681016001600160401b03811182821017156117115761171161165b565b604052818152828201602001891015611728575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b634e487b7160e01b5f52602160045260245ffd5b6003811061177957634e487b7160e01b5f52602160045260245ffd5b9052565b5f6080820190508251825261ffff60208401511660208301526001600160401b03604084015116604083015260608301516117bb606084018261175d565b5092915050565b5f5f5f5f5f5f60a087890312156117d7575f5ffd5b6117e087611562565b955060208701359450604087013561ffff811681146117fd575f5ffd5b935060608701356001600160401b0381168114611818575f5ffd5b925060808701356001600160401b03811115611832575f5ffd5b8701601f81018913611842575f5ffd5b80356001600160401b03811115611857575f5ffd5b896020828401011115611868575f5ffd5b60208201935080925050509295509295509295565b5f5f6040838503121561188e575f5ffd5b61189783611562565b915061160060208401611562565b600181811c908216806118b957607f821691505b6020821081036118d757634e487b7160e01b5f52602260045260245ffd5b50919050565b602081016103c2828461175d565b5f81518060208401855e5f93019283525090919050565b5f61094b61191083866118eb565b846118eb565b5f6001820161193357634e487b7160e01b5f52601160045260245ffd5b5060010190565b6001600160a01b03858116825284166020820152604081018390526080606082018190525f9061196c9083018461150b565b9695505050505050565b5f60208284031215611986575f5ffd5b8151611019816114db565b601f82111561057257805f5260205f20601f840160051c810160208510156119b65750805b601f840160051c820191505b81811015610fa8575f81556001016119c2565b81516001600160401b038111156119ee576119ee61165b565b611a02816119fc84546118a5565b84611991565b6020601f821160018114611a34575f8315611a1d5750848201515b5f19600385901b1c1916600184901b178455610fa8565b5f84815260208120601f198516915b82811015611a635787850151825560209485019460019092019101611a43565b5084821015611a8057868401515f19600387901b60f8161c191681555b50505050600190811b0190555056fea164736f6c634300081c000a",
}

// LoanNoteABI is the input ABI used to generate the binding from.
// Deprecated: Use LoanNoteMetaData.ABI instead.
var LoanNoteABI = LoanNoteMetaData.ABI

// LoanNoteBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LoanNoteMetaData.Bin instead.
var LoanNoteBin = LoanNoteMetaData.Bin

// DeployLoanNote deploys a new Ethereum contract, binding an instance of LoanNote to it.
func DeployLoanNote(auth *bind.TransactOpts, backend bind.ContractBackend, admin common.Address) (common.Address, *types.Transaction, *LoanNote, error) {
	parsed, err := LoanNoteMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LoanNoteBin), backend, admin)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &LoanNote{LoanNoteCaller: LoanNoteCaller{contract: contract}, LoanNoteTransactor: LoanNoteTransactor{contract: contract}, LoanNoteFilterer: LoanNoteFilterer{contract: contract}}, nil
}

// LoanNote is an auto generated Go binding around an Ethereum contract.
type LoanNote struct {
	LoanNoteCaller     // Read-only binding to the contract
	LoanNoteTransactor // Write-only binding to the contract
	LoanNoteFilterer   // Log filterer for contract events
}

// LoanNoteCaller is an auto generated read-only Go binding around an Ethereum contract.
type LoanNoteCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoanNoteTransactor is an auto generated write-only Go binding around an Ethereum contract.
type LoanNoteTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoanNoteFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type LoanNoteFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// LoanNoteSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type LoanNoteSession struct {
	Contract     *LoanNote         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// LoanNoteCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type LoanNoteCallerSession struct {
	Contract *LoanNoteCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// LoanNoteTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type LoanNoteTransactorSession struct {
	Contract     *LoanNoteTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// LoanNoteRaw is an auto generated low-level Go binding around an Ethereum contract.
type LoanNoteRaw struct {
	Contract *LoanNote // Generic contract binding to access the raw methods on
}

// LoanNoteCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type LoanNoteCallerRaw struct {
	Contract *LoanNoteCaller // Generic read-only contract binding to access the raw methods on
}

// LoanNoteTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type LoanNoteTransactorRaw struct {
	Contract *LoanNoteTransactor // Generic write-only contract binding to access the raw methods on
}

// NewLoanNote creates a new instance of LoanNote, bound to a specific deployed contract.
func NewLoanNote(address common.Address, backend bind.ContractBackend) (*LoanNote, error) {
	contract, err := bindLoanNote(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &LoanNote{LoanNoteCaller: LoanNoteCaller{contract: contract}, LoanNoteTransactor: LoanNoteTransactor{contract: contract}, LoanNoteFilterer: LoanNoteFilterer{contract: contract}}, nil
}

// NewLoanNoteCaller creates a new read-only instance of LoanNote, bound to a specific deployed contract.
func NewLoanNoteCaller(address common.Address, caller bind.ContractCaller) (*LoanNoteCaller, error) {
	contract, err := bindLoanNote(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &LoanNoteCaller{contract: contract}, nil
}

// NewLoanNoteTransactor creates a new write-only instance of LoanNote, bound to a specific deployed contract.
func NewLoanNoteTransactor(address common.Address, transactor bind.ContractTransactor) (*LoanNoteTransactor, error) {
	contract, err := bindLoanNote(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &LoanNoteTransactor{contract: contract}, nil
}

// NewLoanNoteFilterer creates a new log filterer instance of LoanNote, bound to a specific deployed contract.
func NewLoanNoteFilterer(address common.Address, filterer bind.ContractFilterer) (*LoanNoteFilterer, error) {
	contract, err := bindLoanNote(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &LoanNoteFilterer{contract: contract}, nil
}

// bindLoanNote binds a generic wrapper to an already deployed contract.
func bindLoanNote(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := LoanNoteMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LoanNote *LoanNoteRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LoanNote.Contract.LoanNoteCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LoanNote *LoanNoteRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoanNote.Contract.LoanNoteTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LoanNote *LoanNoteRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LoanNote.Contract.LoanNoteTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_LoanNote *LoanNoteCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _LoanNote.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_LoanNote *LoanNoteTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoanNote.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_LoanNote *LoanNoteTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _LoanNote.Contract.contract.Transact(opts, method, params...)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_LoanNote *LoanNoteCaller) DEFAULTADMINROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "DEFAULT_ADMIN_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_LoanNote *LoanNoteSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _LoanNote.Contract.DEFAULTADMINROLE(&_LoanNote.CallOpts)
}

// DEFAULTADMINROLE is a free data retrieval call binding the contract method 0xa217fddf.
//
// Solidity: function DEFAULT_ADMIN_ROLE() view returns(bytes32)
func (_LoanNote *LoanNoteCallerSession) DEFAULTADMINROLE() ([32]byte, error) {
	return _LoanNote.Contract.DEFAULTADMINROLE(&_LoanNote.CallOpts)
}

// ORIGINATORROLE is a free data retrieval call binding the contract method 0xa125142b.
//
// Solidity: function ORIGINATOR_ROLE() view returns(bytes32)
func (_LoanNote *LoanNoteCaller) ORIGINATORROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "ORIGINATOR_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// ORIGINATORROLE is a free data retrieval call binding the contract method 0xa125142b.
//
// Solidity: function ORIGINATOR_ROLE() view returns(bytes32)
func (_LoanNote *LoanNoteSession) ORIGINATORROLE() ([32]byte, error) {
	return _LoanNote.Contract.ORIGINATORROLE(&_LoanNote.CallOpts)
}

// ORIGINATORROLE is a free data retrieval call binding the contract method 0xa125142b.
//
// Solidity: function ORIGINATOR_ROLE() view returns(bytes32)
func (_LoanNote *LoanNoteCallerSession) ORIGINATORROLE() ([32]byte, error) {
	return _LoanNote.Contract.ORIGINATORROLE(&_LoanNote.CallOpts)
}

// SERVICERROLE is a free data retrieval call binding the contract method 0x09a01608.
//
// Solidity: function SERVICER_ROLE() view returns(bytes32)
func (_LoanNote *LoanNoteCaller) SERVICERROLE(opts *bind.CallOpts) ([32]byte, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "SERVICER_ROLE")

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// SERVICERROLE is a free data retrieval call binding the contract method 0x09a01608.
//
// Solidity: function SERVICER_ROLE() view returns(bytes32)
func (_LoanNote *LoanNoteSession) SERVICERROLE() ([32]byte, error) {
	return _LoanNote.Contract.SERVICERROLE(&_LoanNote.CallOpts)
}

// SERVICERROLE is a free data retrieval call binding the contract method 0x09a01608.
//
// Solidity: function SERVICER_ROLE() view returns(bytes32)
func (_LoanNote *LoanNoteCallerSession) SERVICERROLE() ([32]byte, error) {
	return _LoanNote.Contract.SERVICERROLE(&_LoanNote.CallOpts)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_LoanNote *LoanNoteCaller) BalanceOf(opts *bind.CallOpts, owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "balanceOf", owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_LoanNote *LoanNoteSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LoanNote.Contract.BalanceOf(&_LoanNote.CallOpts, owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address owner) view returns(uint256)
func (_LoanNote *LoanNoteCallerSession) BalanceOf(owner common.Address) (*big.Int, error) {
	return _LoanNote.Contract.BalanceOf(&_LoanNote.CallOpts, owner)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_LoanNote *LoanNoteCaller) GetApproved(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "getApproved", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_LoanNote *LoanNoteSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _LoanNote.Contract.GetApproved(&_LoanNote.CallOpts, tokenId)
}

// GetApproved is a free data retrieval call binding the contract method 0x081812fc.
//
// Solidity: function getApproved(uint256 tokenId) view returns(address)
func (_LoanNote *LoanNoteCallerSession) GetApproved(tokenId *big.Int) (common.Address, error) {
	return _LoanNote.Contract.GetApproved(&_LoanNote.CallOpts, tokenId)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_LoanNote *LoanNoteCaller) GetRoleAdmin(opts *bind.CallOpts, role [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "getRoleAdmin", role)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_LoanNote *LoanNoteSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _LoanNote.Contract.GetRoleAdmin(&_LoanNote.CallOpts, role)
}

// GetRoleAdmin is a free data retrieval call binding the contract method 0x248a9ca3.
//
// Solidity: function getRoleAdmin(bytes32 role) view returns(bytes32)
func (_LoanNote *LoanNoteCallerSession) GetRoleAdmin(role [32]byte) ([32]byte, error) {
	return _LoanNote.Contract.GetRoleAdmin(&_LoanNote.CallOpts, role)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_LoanNote *LoanNoteCaller) HasRole(opts *bind.CallOpts, role [32]byte, account common.Address) (bool, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "hasRole", role, account)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_LoanNote *LoanNoteSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _LoanNote.Contract.HasRole(&_LoanNote.CallOpts, role, account)
}

// HasRole is a free data retrieval call binding the contract method 0x91d14854.
//
// Solidity: function hasRole(bytes32 role, address account) view returns(bool)
func (_LoanNote *LoanNoteCallerSession) HasRole(role [32]byte, account common.Address) (bool, error) {
	return _LoanNote.Contract.HasRole(&_LoanNote.CallOpts, role, account)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_LoanNote *LoanNoteCaller) IsApprovedForAll(opts *bind.CallOpts, owner common.Address, operator common.Address) (bool, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "isApprovedForAll", owner, operator)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_LoanNote *LoanNoteSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _LoanNote.Contract.IsApprovedForAll(&_LoanNote.CallOpts, owner, operator)
}

// IsApprovedForAll is a free data retrieval call binding the contract method 0xe985e9c5.
//
// Solidity: function isApprovedForAll(address owner, address operator) view returns(bool)
func (_LoanNote *LoanNoteCallerSession) IsApprovedForAll(owner common.Address, operator common.Address) (bool, error) {
	return _LoanNote.Contract.IsApprovedForAll(&_LoanNote.CallOpts, owner, operator)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LoanNote *LoanNoteCaller) Name(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "name")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LoanNote *LoanNoteSession) Name() (string, error) {
	return _LoanNote.Contract.Name(&_LoanNote.CallOpts)
}

// Name is a free data retrieval call binding the contract method 0x06fdde03.
//
// Solidity: function name() view returns(string)
func (_LoanNote *LoanNoteCallerSession) Name() (string, error) {
	return _LoanNote.Contract.Name(&_LoanNote.CallOpts)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_LoanNote *LoanNoteCaller) OwnerOf(opts *bind.CallOpts, tokenId *big.Int) (common.Address, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "ownerOf", tokenId)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_LoanNote *LoanNoteSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _LoanNote.Contract.OwnerOf(&_LoanNote.CallOpts, tokenId)
}

// OwnerOf is a free data retrieval call binding the contract method 0x6352211e.
//
// Solidity: function ownerOf(uint256 tokenId) view returns(address)
func (_LoanNote *LoanNoteCallerSession) OwnerOf(tokenId *big.Int) (common.Address, error) {
	return _LoanNote.Contract.OwnerOf(&_LoanNote.CallOpts, tokenId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_LoanNote *LoanNoteCaller) SupportsInterface(opts *bind.CallOpts, interfaceId [4]byte) (bool, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "supportsInterface", interfaceId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_LoanNote *LoanNoteSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LoanNote.Contract.SupportsInterface(&_LoanNote.CallOpts, interfaceId)
}

// SupportsInterface is a free data retrieval call binding the contract method 0x01ffc9a7.
//
// Solidity: function supportsInterface(bytes4 interfaceId) view returns(bool)
func (_LoanNote *LoanNoteCallerSession) SupportsInterface(interfaceId [4]byte) (bool, error) {
	return _LoanNote.Contract.SupportsInterface(&_LoanNote.CallOpts, interfaceId)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LoanNote *LoanNoteCaller) Symbol(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "symbol")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LoanNote *LoanNoteSession) Symbol() (string, error) {
	return _LoanNote.Contract.Symbol(&_LoanNote.CallOpts)
}

// Symbol is a free data retrieval call binding the contract method 0x95d89b41.
//
// Solidity: function symbol() view returns(string)
func (_LoanNote *LoanNoteCallerSession) Symbol() (string, error) {
	return _LoanNote.Contract.Symbol(&_LoanNote.CallOpts)
}

// Terms is a free data retrieval call binding the contract method 0xc0aa0e8a.
//
// Solidity: function terms(uint256 tokenId) view returns((uint256,uint16,uint64,uint8))
func (_LoanNote *LoanNoteCaller) Terms(opts *bind.CallOpts, tokenId *big.Int) (LoanNoteLoanTerms, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "terms", tokenId)

	if err != nil {
		return *new(LoanNoteLoanTerms), err
	}

	out0 := *abi.ConvertType(out[0], new(LoanNoteLoanTerms)).(*LoanNoteLoanTerms)

	return out0, err

}

// Terms is a free data retrieval call binding the contract method 0xc0aa0e8a.
//
// Solidity: function terms(uint256 tokenId) view returns((uint256,uint16,uint64,uint8))
func (_LoanNote *LoanNoteSession) Terms(tokenId *big.Int) (LoanNoteLoanTerms, error) {
	return _LoanNote.Contract.Terms(&_LoanNote.CallOpts, tokenId)
}

// Terms is a free data retrieval call binding the contract method 0xc0aa0e8a.
//
// Solidity: function terms(uint256 tokenId) view returns((uint256,uint16,uint64,uint8))
func (_LoanNote *LoanNoteCallerSession) Terms(tokenId *big.Int) (LoanNoteLoanTerms, error) {
	return _LoanNote.Contract.Terms(&_LoanNote.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_LoanNote *LoanNoteCaller) TokenURI(opts *bind.CallOpts, tokenId *big.Int) (string, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "tokenURI", tokenId)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_LoanNote *LoanNoteSession) TokenURI(tokenId *big.Int) (string, error) {
	return _LoanNote.Contract.TokenURI(&_LoanNote.CallOpts, tokenId)
}

// TokenURI is a free data retrieval call binding the contract method 0xc87b56dd.
//
// Solidity: function tokenURI(uint256 tokenId) view returns(string)
func (_LoanNote *LoanNoteCallerSession) TokenURI(tokenId *big.Int) (string, error) {
	return _LoanNote.Contract.TokenURI(&_LoanNote.CallOpts, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_LoanNote *LoanNoteTransactor) Approve(opts *bind.TransactOpts, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "approve", to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_LoanNote *LoanNoteSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.Contract.Approve(&_LoanNote.TransactOpts, to, tokenId)
}

// Approve is a paid mutator transaction binding the contract method 0x095ea7b3.
//
// Solidity: function approve(address to, uint256 tokenId) returns()
func (_LoanNote *LoanNoteTransactorSession) Approve(to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.Contract.Approve(&_LoanNote.TransactOpts, to, tokenId)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_LoanNote *LoanNoteTransactor) GrantRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "grantRole", role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_LoanNote *LoanNoteSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoanNote.Contract.GrantRole(&_LoanNote.TransactOpts, role, account)
}

// GrantRole is a paid mutator transaction binding the contract method 0x2f2ff15d.
//
// Solidity: function grantRole(bytes32 role, address account) returns()
func (_LoanNote *LoanNoteTransactorSession) GrantRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoanNote.Contract.GrantRole(&_LoanNote.TransactOpts, role, account)
}

// MarkDefaulted is a paid mutator transaction binding the contract method 0x73216450.
//
// Solidity: function markDefaulted(uint256 tokenId) returns()
func (_LoanNote *LoanNoteTransactor) MarkDefaulted(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "markDefaulted", tokenId)
}

// MarkDefaulted is a paid mutator transaction binding the contract method 0x73216450.
//
// Solidity: function markDefaulted(uint256 tokenId) returns()
func (_LoanNote *LoanNoteSession) MarkDefaulted(tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.Contract.MarkDefaulted(&_LoanNote.TransactOpts, tokenId)
}

// MarkDefaulted is a paid mutator transaction binding the contract method 0x73216450.
//
// Solidity: function markDefaulted(uint256 tokenId) returns()
func (_LoanNote *LoanNoteTransactorSession) MarkDefaulted(tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.Contract.MarkDefaulted(&_LoanNote.TransactOpts, tokenId)
}

// Originate is a paid mutator transaction binding the contract method 0xdaf18fba.
//
// Solidity: function originate(address lender, uint256 principal, uint16 aprBps, uint64 maturity, string tokenURI_) returns(uint256 tokenId)
func (_LoanNote *LoanNoteTransactor) Originate(opts *bind.TransactOpts, lender common.Address, principal *big.Int, aprBps uint16, maturity uint64, tokenURI_ string) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "originate", lender, principal, aprBps, maturity, tokenURI_)
}

// Originate is a paid mutator transaction binding the contract method 0xdaf18fba.
//
// Solidity: function originate(address lender, uint256 principal, uint16 aprBps, uint64 maturity, string tokenURI_) returns(uint256 tokenId)
func (_LoanNote *LoanNoteSession) Originate(lender common.Address, principal *big.Int, aprBps uint16, maturity uint64, tokenURI_ string) (*types.Transaction, error) {
	return _LoanNote.Contract.Originate(&_LoanNote.TransactOpts, lender, principal, aprBps, maturity, tokenURI_)
}

// Originate is a paid mutator transaction binding the contract method 0xdaf18fba.
//
// Solidity: function originate(address lender, uint256 principal, uint16 aprBps, uint64 maturity, string tokenURI_) returns(uint256 tokenId)
func (_LoanNote *LoanNoteTransactorSession) Originate(lender common.Address, principal *big.Int, aprBps uint16, maturity uint64, tokenURI_ string) (*types.Transaction, error) {
	return _LoanNote.Contract.Originate(&_LoanNote.TransactOpts, lender, principal, aprBps, maturity, tokenURI_)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_LoanNote *LoanNoteTransactor) RenounceRole(opts *bind.TransactOpts, role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "renounceRole", role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_LoanNote *LoanNoteSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _LoanNote.Contract.RenounceRole(&_LoanNote.TransactOpts, role, callerConfirmation)
}

// RenounceRole is a paid mutator transaction binding the contract method 0x36568abe.
//
// Solidity: function renounceRole(bytes32 role, address callerConfirmation) returns()
func (_LoanNote *LoanNoteTransactorSession) RenounceRole(role [32]byte, callerConfirmation common.Address) (*types.Transaction, error) {
	return _LoanNote.Contract.RenounceRole(&_LoanNote.TransactOpts, role, callerConfirmation)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_LoanNote *LoanNoteTransactor) RevokeRole(opts *bind.TransactOpts, role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "revokeRole", role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_LoanNote *LoanNoteSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoanNote.Contract.RevokeRole(&_LoanNote.TransactOpts, role, account)
}

// RevokeRole is a paid mutator transaction binding the contract method 0xd547741f.
//
// Solidity: function revokeRole(bytes32 role, address account) returns()
func (_LoanNote *LoanNoteTransactorSession) RevokeRole(role [32]byte, account common.Address) (*types.Transaction, error) {
	return _LoanNote.Contract.RevokeRole(&_LoanNote.TransactOpts, role, account)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_LoanNote *LoanNoteTransactor) SafeTransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "safeTransferFrom", from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_LoanNote *LoanNoteSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.Contract.SafeTransferFrom(&_LoanNote.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom is a paid mutator transaction binding the contract method 0x42842e0e.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId) returns()
func (_LoanNote *LoanNoteTransactorSession) SafeTransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.Contract.SafeTransferFrom(&_LoanNote.TransactOpts, from, to, tokenId)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_LoanNote *LoanNoteTransactor) SafeTransferFrom0(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "safeTransferFrom0", from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_LoanNote *LoanNoteSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _LoanNote.Contract.SafeTransferFrom0(&_LoanNote.TransactOpts, from, to, tokenId, data)
}

// SafeTransferFrom0 is a paid mutator transaction binding the contract method 0xb88d4fde.
//
// Solidity: function safeTransferFrom(address from, address to, uint256 tokenId, bytes data) returns()
func (_LoanNote *LoanNoteTransactorSession) SafeTransferFrom0(from common.Address, to common.Address, tokenId *big.Int, data []byte) (*types.Transaction, error) {
	return _LoanNote.Contract.SafeTransferFrom0(&_LoanNote.TransactOpts, from, to, tokenId, data)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_LoanNote *LoanNoteTransactor) SetApprovalForAll(opts *bind.TransactOpts, operator common.Address, approved bool) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "setApprovalForAll", operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_LoanNote *LoanNoteSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _LoanNote.Contract.SetApprovalForAll(&_LoanNote.TransactOpts, operator, approved)
}

// SetApprovalForAll is a paid mutator transaction binding the contract method 0xa22cb465.
//
// Solidity: function setApprovalForAll(address operator, bool approved) returns()
func (_LoanNote *LoanNoteTransactorSession) SetApprovalForAll(operator common.Address, approved bool) (*types.Transaction, error) {
	return _LoanNote.Contract.SetApprovalForAll(&_LoanNote.TransactOpts, operator, approved)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 tokenId) returns()
func (_LoanNote *LoanNoteTransactor) Settle(opts *bind.TransactOpts, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "settle", tokenId)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 tokenId) returns()
func (_LoanNote *LoanNoteSession) Settle(tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.Contract.Settle(&_LoanNote.TransactOpts, tokenId)
}

// Settle is a paid mutator transaction binding the contract method 0x8df82800.
//
// Solidity: function settle(uint256 tokenId) returns()
func (_LoanNote *LoanNoteTransactorSession) Settle(tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.Contract.Settle(&_LoanNote.TransactOpts, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_LoanNote *LoanNoteTransactor) TransferFrom(opts *bind.TransactOpts, from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "transferFrom", from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_LoanNote *LoanNoteSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.Contract.TransferFrom(&_LoanNote.TransactOpts, from, to, tokenId)
}

// TransferFrom is a paid mutator transaction binding the contract method 0x23b872dd.
//
// Solidity: function transferFrom(address from, address to, uint256 tokenId) returns()
func (_LoanNote *LoanNoteTransactorSession) TransferFrom(from common.Address, to common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _LoanNote.Contract.TransferFrom(&_LoanNote.TransactOpts, from, to, tokenId)
}

// LoanNoteApprovalIterator is returned from FilterApproval and is used to iterate over the raw logs and unpacked data for Approval events raised by the LoanNote contract.
type LoanNoteApprovalIterator struct {
	Event *LoanNoteApproval // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteApprovalIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteApproval)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteApproval)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteApprovalIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteApprovalIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteApproval represents a Approval event raised by the LoanNote contract.
type LoanNoteApproval struct {
	Owner    common.Address
	Approved common.Address
	TokenId  *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApproval is a free log retrieval operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_LoanNote *LoanNoteFilterer) FilterApproval(opts *bind.FilterOpts, owner []common.Address, approved []common.Address, tokenId []*big.Int) (*LoanNoteApprovalIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &LoanNoteApprovalIterator{contract: _LoanNote.contract, event: "Approval", logs: logs, sub: sub}, nil
}

// WatchApproval is a free log subscription operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_LoanNote *LoanNoteFilterer) WatchApproval(opts *bind.WatchOpts, sink chan<- *LoanNoteApproval, owner []common.Address, approved []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var approvedRule []interface{}
	for _, approvedItem := range approved {
		approvedRule = append(approvedRule, approvedItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "Approval", ownerRule, approvedRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteApproval)
				if err := _LoanNote.contract.UnpackLog(event, "Approval", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApproval is a log parse operation binding the contract event 0x8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b925.
//
// Solidity: event Approval(address indexed owner, address indexed approved, uint256 indexed tokenId)
func (_LoanNote *LoanNoteFilterer) ParseApproval(log types.Log) (*LoanNoteApproval, error) {
	event := new(LoanNoteApproval)
	if err := _LoanNote.contract.UnpackLog(event, "Approval", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoanNoteApprovalForAllIterator is returned from FilterApprovalForAll and is used to iterate over the raw logs and unpacked data for ApprovalForAll events raised by the LoanNote contract.
type LoanNoteApprovalForAllIterator struct {
	Event *LoanNoteApprovalForAll // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteApprovalForAllIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteApprovalForAll)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteApprovalForAll)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteApprovalForAllIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteApprovalForAllIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteApprovalForAll represents a ApprovalForAll event raised by the LoanNote contract.
type LoanNoteApprovalForAll struct {
	Owner    common.Address
	Operator common.Address
	Approved bool
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterApprovalForAll is a free log retrieval operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_LoanNote *LoanNoteFilterer) FilterApprovalForAll(opts *bind.FilterOpts, owner []common.Address, operator []common.Address) (*LoanNoteApprovalForAllIterator, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return &LoanNoteApprovalForAllIterator{contract: _LoanNote.contract, event: "ApprovalForAll", logs: logs, sub: sub}, nil
}

// WatchApprovalForAll is a free log subscription operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_LoanNote *LoanNoteFilterer) WatchApprovalForAll(opts *bind.WatchOpts, sink chan<- *LoanNoteApprovalForAll, owner []common.Address, operator []common.Address) (event.Subscription, error) {

	var ownerRule []interface{}
	for _, ownerItem := range owner {
		ownerRule = append(ownerRule, ownerItem)
	}
	var operatorRule []interface{}
	for _, operatorItem := range operator {
		operatorRule = append(operatorRule, operatorItem)
	}

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "ApprovalForAll", ownerRule, operatorRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteApprovalForAll)
				if err := _LoanNote.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseApprovalForAll is a log parse operation binding the contract event 0x17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31.
//
// Solidity: event ApprovalForAll(address indexed owner, address indexed operator, bool approved)
func (_LoanNote *LoanNoteFilterer) ParseApprovalForAll(log types.Log) (*LoanNoteApprovalForAll, error) {
	event := new(LoanNoteApprovalForAll)
	if err := _LoanNote.contract.UnpackLog(event, "ApprovalForAll", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoanNoteBatchMetadataUpdateIterator is returned from FilterBatchMetadataUpdate and is used to iterate over the raw logs and unpacked data for BatchMetadataUpdate events raised by the LoanNote contract.
type LoanNoteBatchMetadataUpdateIterator struct {
	Event *LoanNoteBatchMetadataUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteBatchMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteBatchMetadataUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteBatchMetadataUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteBatchMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteBatchMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteBatchMetadataUpdate represents a BatchMetadataUpdate event raised by the LoanNote contract.
type LoanNoteBatchMetadataUpdate struct {
	FromTokenId *big.Int
	ToTokenId   *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterBatchMetadataUpdate is a free log retrieval operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_LoanNote *LoanNoteFilterer) FilterBatchMetadataUpdate(opts *bind.FilterOpts) (*LoanNoteBatchMetadataUpdateIterator, error) {

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &LoanNoteBatchMetadataUpdateIterator{contract: _LoanNote.contract, event: "BatchMetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchBatchMetadataUpdate is a free log subscription operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_LoanNote *LoanNoteFilterer) WatchBatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *LoanNoteBatchMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "BatchMetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteBatchMetadataUpdate)
				if err := _LoanNote.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseBatchMetadataUpdate is a log parse operation binding the contract event 0x6bd5c950a8d8df17f772f5af37cb3655737899cbf903264b9795592da439661c.
//
// Solidity: event BatchMetadataUpdate(uint256 _fromTokenId, uint256 _toTokenId)
func (_LoanNote *LoanNoteFilterer) ParseBatchMetadataUpdate(log types.Log) (*LoanNoteBatchMetadataUpdate, error) {
	event := new(LoanNoteBatchMetadataUpdate)
	if err := _LoanNote.contract.UnpackLog(event, "BatchMetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoanNoteLoanOriginatedIterator is returned from FilterLoanOriginated and is used to iterate over the raw logs and unpacked data for LoanOriginated events raised by the LoanNote contract.
type LoanNoteLoanOriginatedIterator struct {
	Event *LoanNoteLoanOriginated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteLoanOriginatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteLoanOriginated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteLoanOriginated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteLoanOriginatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteLoanOriginatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteLoanOriginated represents a LoanOriginated event raised by the LoanNote contract.
type LoanNoteLoanOriginated struct {
	TokenId   *big.Int
	Lender    common.Address
	Principal *big.Int
	AprBps    uint16
	Maturity  uint64
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterLoanOriginated is a free log retrieval operation binding the contract event 0x3ade124d39b4fb6fb328df4ccd68907e63c76c84fe020c179e2195101b464400.
//
// Solidity: event LoanOriginated(uint256 indexed tokenId, address indexed lender, uint256 principal, uint16 aprBps, uint64 maturity)
func (_LoanNote *LoanNoteFilterer) FilterLoanOriginated(opts *bind.FilterOpts, tokenId []*big.Int, lender []common.Address) (*LoanNoteLoanOriginatedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var lenderRule []interface{}
	for _, lenderItem := range lender {
		lenderRule = append(lenderRule, lenderItem)
	}

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "LoanOriginated", tokenIdRule, lenderRule)
	if err != nil {
		return nil, err
	}
	return &LoanNoteLoanOriginatedIterator{contract: _LoanNote.contract, event: "LoanOriginated", logs: logs, sub: sub}, nil
}

// WatchLoanOriginated is a free log subscription operation binding the contract event 0x3ade124d39b4fb6fb328df4ccd68907e63c76c84fe020c179e2195101b464400.
//
// Solidity: event LoanOriginated(uint256 indexed tokenId, address indexed lender, uint256 principal, uint16 aprBps, uint64 maturity)
func (_LoanNote *LoanNoteFilterer) WatchLoanOriginated(opts *bind.WatchOpts, sink chan<- *LoanNoteLoanOriginated, tokenId []*big.Int, lender []common.Address) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}
	var lenderRule []interface{}
	for _, lenderItem := range lender {
		lenderRule = append(lenderRule, lenderItem)
	}

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "LoanOriginated", tokenIdRule, lenderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteLoanOriginated)
				if err := _LoanNote.contract.UnpackLog(event, "LoanOriginated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLoanOriginated is a log parse operation binding the contract event 0x3ade124d39b4fb6fb328df4ccd68907e63c76c84fe020c179e2195101b464400.
//
// Solidity: event LoanOriginated(uint256 indexed tokenId, address indexed lender, uint256 principal, uint16 aprBps, uint64 maturity)
func (_LoanNote *LoanNoteFilterer) ParseLoanOriginated(log types.Log) (*LoanNoteLoanOriginated, error) {
	event := new(LoanNoteLoanOriginated)
	if err := _LoanNote.contract.UnpackLog(event, "LoanOriginated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoanNoteLoanSettledIterator is returned from FilterLoanSettled and is used to iterate over the raw logs and unpacked data for LoanSettled events raised by the LoanNote contract.
type LoanNoteLoanSettledIterator struct {
	Event *LoanNoteLoanSettled // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteLoanSettledIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteLoanSettled)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteLoanSettled)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteLoanSettledIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteLoanSettledIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteLoanSettled represents a LoanSettled event raised by the LoanNote contract.
type LoanNoteLoanSettled struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLoanSettled is a free log retrieval operation binding the contract event 0xf2257195ce09a74d98ff0578eba9d4573bc34bb4b21473e363535fe3485e0efa.
//
// Solidity: event LoanSettled(uint256 indexed tokenId)
func (_LoanNote *LoanNoteFilterer) FilterLoanSettled(opts *bind.FilterOpts, tokenId []*big.Int) (*LoanNoteLoanSettledIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "LoanSettled", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &LoanNoteLoanSettledIterator{contract: _LoanNote.contract, event: "LoanSettled", logs: logs, sub: sub}, nil
}

// WatchLoanSettled is a free log subscription operation binding the contract event 0xf2257195ce09a74d98ff0578eba9d4573bc34bb4b21473e363535fe3485e0efa.
//
// Solidity: event LoanSettled(uint256 indexed tokenId)
func (_LoanNote *LoanNoteFilterer) WatchLoanSettled(opts *bind.WatchOpts, sink chan<- *LoanNoteLoanSettled, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "LoanSettled", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteLoanSettled)
				if err := _LoanNote.contract.UnpackLog(event, "LoanSettled", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLoanSettled is a log parse operation binding the contract event 0xf2257195ce09a74d98ff0578eba9d4573bc34bb4b21473e363535fe3485e0efa.
//
// Solidity: event LoanSettled(uint256 indexed tokenId)
func (_LoanNote *LoanNoteFilterer) ParseLoanSettled(log types.Log) (*LoanNoteLoanSettled, error) {
	event := new(LoanNoteLoanSettled)
	if err := _LoanNote.contract.UnpackLog(event, "LoanSettled", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoanNoteLoanStatusChangedIterator is returned from FilterLoanStatusChanged and is used to iterate over the raw logs and unpacked data for LoanStatusChanged events raised by the LoanNote contract.
type LoanNoteLoanStatusChangedIterator struct {
	Event *LoanNoteLoanStatusChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteLoanStatusChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteLoanStatusChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteLoanStatusChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteLoanStatusChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteLoanStatusChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteLoanStatusChanged represents a LoanStatusChanged event raised by the LoanNote contract.
type LoanNoteLoanStatusChanged struct {
	TokenId *big.Int
	Status  uint8
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterLoanStatusChanged is a free log retrieval operation binding the contract event 0x9d0b8f6161220422fcfcf3cbe3b12d5148060bea52d7d74395488cae75d2e46f.
//
// Solidity: event LoanStatusChanged(uint256 indexed tokenId, uint8 status)
func (_LoanNote *LoanNoteFilterer) FilterLoanStatusChanged(opts *bind.FilterOpts, tokenId []*big.Int) (*LoanNoteLoanStatusChangedIterator, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "LoanStatusChanged", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &LoanNoteLoanStatusChangedIterator{contract: _LoanNote.contract, event: "LoanStatusChanged", logs: logs, sub: sub}, nil
}

// WatchLoanStatusChanged is a free log subscription operation binding the contract event 0x9d0b8f6161220422fcfcf3cbe3b12d5148060bea52d7d74395488cae75d2e46f.
//
// Solidity: event LoanStatusChanged(uint256 indexed tokenId, uint8 status)
func (_LoanNote *LoanNoteFilterer) WatchLoanStatusChanged(opts *bind.WatchOpts, sink chan<- *LoanNoteLoanStatusChanged, tokenId []*big.Int) (event.Subscription, error) {

	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "LoanStatusChanged", tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteLoanStatusChanged)
				if err := _LoanNote.contract.UnpackLog(event, "LoanStatusChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseLoanStatusChanged is a log parse operation binding the contract event 0x9d0b8f6161220422fcfcf3cbe3b12d5148060bea52d7d74395488cae75d2e46f.
//
// Solidity: event LoanStatusChanged(uint256 indexed tokenId, uint8 status)
func (_LoanNote *LoanNoteFilterer) ParseLoanStatusChanged(log types.Log) (*LoanNoteLoanStatusChanged, error) {
	event := new(LoanNoteLoanStatusChanged)
	if err := _LoanNote.contract.UnpackLog(event, "LoanStatusChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoanNoteMetadataUpdateIterator is returned from FilterMetadataUpdate and is used to iterate over the raw logs and unpacked data for MetadataUpdate events raised by the LoanNote contract.
type LoanNoteMetadataUpdateIterator struct {
	Event *LoanNoteMetadataUpdate // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteMetadataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteMetadataUpdate)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteMetadataUpdate)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteMetadataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteMetadataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteMetadataUpdate represents a MetadataUpdate event raised by the LoanNote contract.
type LoanNoteMetadataUpdate struct {
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterMetadataUpdate is a free log retrieval operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_LoanNote *LoanNoteFilterer) FilterMetadataUpdate(opts *bind.FilterOpts) (*LoanNoteMetadataUpdateIterator, error) {

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return &LoanNoteMetadataUpdateIterator{contract: _LoanNote.contract, event: "MetadataUpdate", logs: logs, sub: sub}, nil
}

// WatchMetadataUpdate is a free log subscription operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_LoanNote *LoanNoteFilterer) WatchMetadataUpdate(opts *bind.WatchOpts, sink chan<- *LoanNoteMetadataUpdate) (event.Subscription, error) {

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "MetadataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteMetadataUpdate)
				if err := _LoanNote.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMetadataUpdate is a log parse operation binding the contract event 0xf8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce7.
//
// Solidity: event MetadataUpdate(uint256 _tokenId)
func (_LoanNote *LoanNoteFilterer) ParseMetadataUpdate(log types.Log) (*LoanNoteMetadataUpdate, error) {
	event := new(LoanNoteMetadataUpdate)
	if err := _LoanNote.contract.UnpackLog(event, "MetadataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoanNoteRoleAdminChangedIterator is returned from FilterRoleAdminChanged and is used to iterate over the raw logs and unpacked data for RoleAdminChanged events raised by the LoanNote contract.
type LoanNoteRoleAdminChangedIterator struct {
	Event *LoanNoteRoleAdminChanged // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteRoleAdminChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteRoleAdminChanged)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteRoleAdminChanged)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteRoleAdminChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteRoleAdminChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteRoleAdminChanged represents a RoleAdminChanged event raised by the LoanNote contract.
type LoanNoteRoleAdminChanged struct {
	Role              [32]byte
	PreviousAdminRole [32]byte
	NewAdminRole      [32]byte
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterRoleAdminChanged is a free log retrieval operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_LoanNote *LoanNoteFilterer) FilterRoleAdminChanged(opts *bind.FilterOpts, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (*LoanNoteRoleAdminChangedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return &LoanNoteRoleAdminChangedIterator{contract: _LoanNote.contract, event: "RoleAdminChanged", logs: logs, sub: sub}, nil
}

// WatchRoleAdminChanged is a free log subscription operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_LoanNote *LoanNoteFilterer) WatchRoleAdminChanged(opts *bind.WatchOpts, sink chan<- *LoanNoteRoleAdminChanged, role [][32]byte, previousAdminRole [][32]byte, newAdminRole [][32]byte) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var previousAdminRoleRule []interface{}
	for _, previousAdminRoleItem := range previousAdminRole {
		previousAdminRoleRule = append(previousAdminRoleRule, previousAdminRoleItem)
	}
	var newAdminRoleRule []interface{}
	for _, newAdminRoleItem := range newAdminRole {
		newAdminRoleRule = append(newAdminRoleRule, newAdminRoleItem)
	}

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "RoleAdminChanged", roleRule, previousAdminRoleRule, newAdminRoleRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteRoleAdminChanged)
				if err := _LoanNote.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleAdminChanged is a log parse operation binding the contract event 0xbd79b86ffe0ab8e8776151514217cd7cacd52c909f66475c3af44e129f0b00ff.
//
// Solidity: event RoleAdminChanged(bytes32 indexed role, bytes32 indexed previousAdminRole, bytes32 indexed newAdminRole)
func (_LoanNote *LoanNoteFilterer) ParseRoleAdminChanged(log types.Log) (*LoanNoteRoleAdminChanged, error) {
	event := new(LoanNoteRoleAdminChanged)
	if err := _LoanNote.contract.UnpackLog(event, "RoleAdminChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoanNoteRoleGrantedIterator is returned from FilterRoleGranted and is used to iterate over the raw logs and unpacked data for RoleGranted events raised by the LoanNote contract.
type LoanNoteRoleGrantedIterator struct {
	Event *LoanNoteRoleGranted // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteRoleGrantedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteRoleGranted)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteRoleGranted)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteRoleGrantedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteRoleGrantedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteRoleGranted represents a RoleGranted event raised by the LoanNote contract.
type LoanNoteRoleGranted struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleGranted is a free log retrieval operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_LoanNote *LoanNoteFilterer) FilterRoleGranted(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*LoanNoteRoleGrantedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &LoanNoteRoleGrantedIterator{contract: _LoanNote.contract, event: "RoleGranted", logs: logs, sub: sub}, nil
}

// WatchRoleGranted is a free log subscription operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_LoanNote *LoanNoteFilterer) WatchRoleGranted(opts *bind.WatchOpts, sink chan<- *LoanNoteRoleGranted, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "RoleGranted", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteRoleGranted)
				if err := _LoanNote.contract.UnpackLog(event, "RoleGranted", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleGranted is a log parse operation binding the contract event 0x2f8788117e7eff1d82e926ec794901d17c78024a50270940304540a733656f0d.
//
// Solidity: event RoleGranted(bytes32 indexed role, address indexed account, address indexed sender)
func (_LoanNote *LoanNoteFilterer) ParseRoleGranted(log types.Log) (*LoanNoteRoleGranted, error) {
	event := new(LoanNoteRoleGranted)
	if err := _LoanNote.contract.UnpackLog(event, "RoleGranted", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoanNoteRoleRevokedIterator is returned from FilterRoleRevoked and is used to iterate over the raw logs and unpacked data for RoleRevoked events raised by the LoanNote contract.
type LoanNoteRoleRevokedIterator struct {
	Event *LoanNoteRoleRevoked // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteRoleRevokedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteRoleRevoked)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteRoleRevoked)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteRoleRevokedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteRoleRevokedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteRoleRevoked represents a RoleRevoked event raised by the LoanNote contract.
type LoanNoteRoleRevoked struct {
	Role    [32]byte
	Account common.Address
	Sender  common.Address
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterRoleRevoked is a free log retrieval operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_LoanNote *LoanNoteFilterer) FilterRoleRevoked(opts *bind.FilterOpts, role [][32]byte, account []common.Address, sender []common.Address) (*LoanNoteRoleRevokedIterator, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return &LoanNoteRoleRevokedIterator{contract: _LoanNote.contract, event: "RoleRevoked", logs: logs, sub: sub}, nil
}

// WatchRoleRevoked is a free log subscription operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_LoanNote *LoanNoteFilterer) WatchRoleRevoked(opts *bind.WatchOpts, sink chan<- *LoanNoteRoleRevoked, role [][32]byte, account []common.Address, sender []common.Address) (event.Subscription, error) {

	var roleRule []interface{}
	for _, roleItem := range role {
		roleRule = append(roleRule, roleItem)
	}
	var accountRule []interface{}
	for _, accountItem := range account {
		accountRule = append(accountRule, accountItem)
	}
	var senderRule []interface{}
	for _, senderItem := range sender {
		senderRule = append(senderRule, senderItem)
	}

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "RoleRevoked", roleRule, accountRule, senderRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteRoleRevoked)
				if err := _LoanNote.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseRoleRevoked is a log parse operation binding the contract event 0xf6391f5c32d9c69d2a47ea670b442974b53935d1edc7fd64eb21e047a839171b.
//
// Solidity: event RoleRevoked(bytes32 indexed role, address indexed account, address indexed sender)
func (_LoanNote *LoanNoteFilterer) ParseRoleRevoked(log types.Log) (*LoanNoteRoleRevoked, error) {
	event := new(LoanNoteRoleRevoked)
	if err := _LoanNote.contract.UnpackLog(event, "RoleRevoked", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// LoanNoteTransferIterator is returned from FilterTransfer and is used to iterate over the raw logs and unpacked data for Transfer events raised by the LoanNote contract.
type LoanNoteTransferIterator struct {
	Event *LoanNoteTransfer // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *LoanNoteTransferIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteTransfer)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(LoanNoteTransfer)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *LoanNoteTransferIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteTransferIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteTransfer represents a Transfer event raised by the LoanNote contract.
type LoanNoteTransfer struct {
	From    common.Address
	To      common.Address
	TokenId *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterTransfer is a free log retrieval operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_LoanNote *LoanNoteFilterer) FilterTransfer(opts *bind.FilterOpts, from []common.Address, to []common.Address, tokenId []*big.Int) (*LoanNoteTransferIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return &LoanNoteTransferIterator{contract: _LoanNote.contract, event: "Transfer", logs: logs, sub: sub}, nil
}

// WatchTransfer is a free log subscription operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_LoanNote *LoanNoteFilterer) WatchTransfer(opts *bind.WatchOpts, sink chan<- *LoanNoteTransfer, from []common.Address, to []common.Address, tokenId []*big.Int) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}
	var toRule []interface{}
	for _, toItem := range to {
		toRule = append(toRule, toItem)
	}
	var tokenIdRule []interface{}
	for _, tokenIdItem := range tokenId {
		tokenIdRule = append(tokenIdRule, tokenIdItem)
	}

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "Transfer", fromRule, toRule, tokenIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteTransfer)
				if err := _LoanNote.contract.UnpackLog(event, "Transfer", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseTransfer is a log parse operation binding the contract event 0xddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef.
//
// Solidity: event Transfer(address indexed from, address indexed to, uint256 indexed tokenId)
func (_LoanNote *LoanNoteFilterer) ParseTransfer(log types.Log) (*LoanNoteTransfer, error) {
	event := new(LoanNoteTransfer)
	if err := _LoanNote.contract.UnpackLog(event, "Transfer", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
