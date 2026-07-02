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
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"initialOwner\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721IncorrectOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721InsufficientApproval\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"approver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidApprover\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOperator\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"ERC721InvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"receiver\",\"type\":\"address\"}],\"name\":\"ERC721InvalidReceiver\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"sender\",\"type\":\"address\"}],\"name\":\"ERC721InvalidSender\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ERC721NonexistentToken\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"approved\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Approval\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"ApprovalForAll\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_fromTokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_toTokenId\",\"type\":\"uint256\"}],\"name\":\"BatchMetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"lender\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint16\",\"name\":\"aprBps\",\"type\":\"uint16\"},{\"indexed\":false,\"internalType\":\"uint64\",\"name\":\"maturity\",\"type\":\"uint64\"}],\"name\":\"LoanOriginated\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"LoanSettled\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"enumLoanNote.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"name\":\"LoanStatusChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"_tokenId\",\"type\":\"uint256\"}],\"name\":\"MetadataUpdate\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"Transfer\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"}],\"name\":\"adminTransfer\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"approve\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"getApproved\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"}],\"name\":\"isApprovedForAll\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"markDefaulted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"name\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"lender\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"aprBps\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"maturity\",\"type\":\"uint64\"},{\"internalType\":\"string\",\"name\":\"tokenURI_\",\"type\":\"string\"}],\"name\":\"originate\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ownerOf\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"safeTransferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"operator\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"approved\",\"type\":\"bool\"}],\"name\":\"setApprovalForAll\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"settle\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes4\",\"name\":\"interfaceId\",\"type\":\"bytes4\"}],\"name\":\"supportsInterface\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"symbol\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"terms\",\"outputs\":[{\"components\":[{\"internalType\":\"uint256\",\"name\":\"principal\",\"type\":\"uint256\"},{\"internalType\":\"uint16\",\"name\":\"aprBps\",\"type\":\"uint16\"},{\"internalType\":\"uint64\",\"name\":\"maturity\",\"type\":\"uint64\"},{\"internalType\":\"enumLoanNote.Status\",\"name\":\"status\",\"type\":\"uint8\"}],\"internalType\":\"structLoanNote.LoanTerms\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"tokenURI\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"transferFrom\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
	Bin: "0x608060405234801561000f575f5ffd5b50604051611bec380380611bec83398101604081905261002e91610119565b80604051806040016040528060088152602001674c6f616e4e6f746560c01b815250604051806040016040528060048152602001632627a0a760e11b815250815f908161007b91906101de565b50600161008882826101de565b5050506001600160a01b0381166100b857604051631e4fbdf760e01b81525f600482015260240160405180910390fd5b6100c1816100c8565b5050610298565b600780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b5f60208284031215610129575f5ffd5b81516001600160a01b038116811461013f575f5ffd5b9392505050565b634e487b7160e01b5f52604160045260245ffd5b600181811c9082168061016e57607f821691505b60208210810361018c57634e487b7160e01b5f52602260045260245ffd5b50919050565b601f8211156101d957805f5260205f20601f840160051c810160208510156101b75750805b601f840160051c820191505b818110156101d6575f81556001016101c3565b50505b505050565b81516001600160401b038111156101f7576101f7610146565b61020b81610205845461015a565b84610192565b6020601f82116001811461023d575f83156102265750848201515b5f19600385901b1c1916600184901b1784556101d6565b5f84815260208120601f198516915b8281101561026c578785015182556020948501946001909201910161024c565b508482101561028957868401515f19600387901b60f8161c191681555b50505050600190811b01905550565b611947806102a55f395ff3fe608060405234801561000f575f5ffd5b5060043610610132575f3560e01c806373216450116100b4578063b88d4fde11610079578063b88d4fde1461027a578063c0aa0e8a1461028d578063c87b56dd146102ad578063daf18fba146102c0578063e985e9c5146102d3578063f2fde38b146102e6575f5ffd5b806373216450146102285780638da5cb5b1461023b5780638df828001461024c57806395d89b411461025f578063a22cb46514610267575f5ffd5b806323b872dd116100fa57806323b872dd146101c657806342842e0e146101d95780636352211e146101ec57806370a08231146101ff578063715018a614610220575f5ffd5b806301ffc9a7146101365780630483fc5e1461015e57806306fdde0314610173578063081812fc14610188578063095ea7b3146101b3575b5f5ffd5b610149610144366004611372565b6102f9565b60405190151581526020015b60405180910390f35b61017161016c3660046113a8565b610323565b005b61017b610347565b6040516101559190611400565b61019b610196366004611412565b6103d6565b6040516001600160a01b039091168152602001610155565b6101716101c1366004611429565b6103fd565b6101716101d4366004611451565b61040c565b6101716101e7366004611451565b61049a565b61019b6101fa366004611412565b6104b4565b61021261020d36600461148b565b6104be565b604051908152602001610155565b610171610503565b610171610236366004611412565b610516565b6007546001600160a01b031661019b565b61017161025a366004611412565b610587565b61017b61062c565b6101716102753660046114a4565b61063b565b6101716102883660046114f1565b610646565b6102a061029b366004611412565b61065e565b60405161015591906115ff565b61017b6102bb366004611412565b61070e565b6102126102ce366004611644565b610819565b6101496102e13660046116ff565b61099f565b6101716102f436600461148b565b6109cc565b5f6001600160e01b03198216632483248360e11b148061031d575061031d82610a09565b92915050565b61032b610a58565b5f610335836104b4565b9050610342818385610a85565b505050565b60605f805461035590611727565b80601f016020809104026020016040519081016040528092919081815260200182805461038190611727565b80156103cc5780601f106103a3576101008083540402835291602001916103cc565b820191905f5260205f20905b8154815290600101906020018083116103af57829003601f168201915b5050505050905090565b5f6103e082610a9f565b505f828152600460205260409020546001600160a01b031661031d565b610408828233610ad7565b5050565b6001600160a01b03821661043a57604051633250574960e11b81525f60048201526024015b60405180910390fd5b5f610446838333610ae4565b9050836001600160a01b0316816001600160a01b031614610494576040516364283d7b60e01b81526001600160a01b0380861660048301526024820184905282166044820152606401610431565b50505050565b61034283838360405180602001604052805f815250610646565b5f61031d82610a9f565b5f6001600160a01b0382166104e8576040516322718ad960e21b81525f6004820152602401610431565b506001600160a01b03165f9081526003602052604090205490565b61050b610a58565b6105145f610bd6565b565b61051e610a58565b61052781610a9f565b505f8181526009602052604090819020600101805460ff60501b1916600160511b1790555181907f9d0b8f6161220422fcfcf3cbe3b12d5148060bea52d7d74395488cae75d2e46f9061057c9060029061175f565b60405180910390a250565b61058f610a58565b61059881610a9f565b50807f9d0b8f6161220422fcfcf3cbe3b12d5148060bea52d7d74395488cae75d2e46f60016040516105ca919061175f565b60405180910390a260405181907ff2257195ce09a74d98ff0578eba9d4573bc34bb4b21473e363535fe3485e0efa905f90a261060581610c27565b5f90815260096020526040812090815560010180546affffffffffffffffffffff19169055565b60606001805461035590611727565b610408338383610c5f565b61065184848461040c565b6104943385858585610cfd565b610685604080516080810182525f8082526020820181905291810182905290606082015290565b61068e82610a9f565b505f82815260096020908152604091829020825160808101845281548152600182015461ffff8116938201939093526201000083046001600160401b031693810193909352906060830190600160501b900460ff1660028111156106f4576106f46115cb565b6002811115610705576107056115cb565b90525092915050565b606061071982610a9f565b505f828152600660205260408120805461073290611727565b80601f016020809104026020016040519081016040528092919081815260200182805461075e90611727565b80156107a95780601f10610780576101008083540402835291602001916107a9565b820191905f5260205f20905b81548152906001019060200180831161078c57829003601f168201915b505050505090505f6107c560408051602081019091525f815290565b905080515f036107d6575092915050565b8151156108085780826040516020016107f0929190611784565b60405160208183030381529060405292505050919050565b61081184610e25565b949350505050565b5f610822610a58565b60088054905f61083183611798565b91905055905060405180608001604052808781526020018661ffff168152602001856001600160401b031681526020015f6002811115610873576108736115cb565b90525f828152600960209081526040918290208351815590830151600182018054938501516001600160401b0316620100000269ffffffffffffffffffff1990941661ffff9092169190911792909217808355606084015191929060ff60501b1916600160501b8360028111156108ec576108ec6115cb565b02179055509050506108fe8782610e96565b61093d8184848080601f0160208091040260200160405190810160405280939291908181526020018383808284375f92019190915250610eaf92505050565b6040805187815261ffff871660208201526001600160401b0386168183015290516001600160a01b0389169183917f3ade124d39b4fb6fb328df4ccd68907e63c76c84fe020c179e2195101b4644009181900360600190a39695505050505050565b6001600160a01b039182165f90815260056020908152604080832093909416825291909152205460ff1690565b6109d4610a58565b6001600160a01b0381166109fd57604051631e4fbdf760e01b81525f6004820152602401610431565b610a0681610bd6565b50565b5f6001600160e01b031982166380ac58cd60e01b1480610a3957506001600160e01b03198216635b5e139f60e01b145b8061031d57506301ffc9a760e01b6001600160e01b031983161461031d565b6007546001600160a01b031633146105145760405163118cdaa760e01b8152336004820152602401610431565b61034283838360405180602001604052805f815250610efe565b5f818152600260205260408120546001600160a01b03168061031d57604051637e27328960e01b815260048101849052602401610431565b6103428383836001610f09565b5f828152600260205260408120546001600160a01b0390811690831615610b1057610b1081848661100d565b6001600160a01b03811615610b4a57610b2b5f855f5f610f09565b6001600160a01b0381165f90815260036020526040902080545f190190555b6001600160a01b03851615610b78576001600160a01b0385165f908152600360205260409020805460010190555b5f8481526002602052604080822080546001600160a01b0319166001600160a01b0389811691821790925591518793918516917fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef91a4949350505050565b600780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0905f90a35050565b5f610c335f835f610ae4565b90506001600160a01b03811661040857604051637e27328960e01b815260048101839052602401610431565b6001600160a01b038216610c9157604051630b61174360e31b81526001600160a01b0383166004820152602401610431565b6001600160a01b038381165f81815260056020908152604080832094871680845294825291829020805460ff191686151590811790915591519182527f17307eab39ab6107e8899845ad3d59bd9653f200f220920489ca2b5937696c31910160405180910390a3505050565b6001600160a01b0383163b15610e1e57604051630a85bd0160e11b81526001600160a01b0384169063150b7a0290610d3f9088908890879087906004016117bc565b6020604051808303815f875af1925050508015610d79575060408051601f3d908101601f19168201909252610d76918101906117f8565b60015b610de0573d808015610da6576040519150601f19603f3d011682016040523d82523d5f602084013e610dab565b606091505b5080515f03610dd857604051633250574960e11b81526001600160a01b0385166004820152602401610431565b805181602001fd5b6001600160e01b03198116630a85bd0160e11b14610e1c57604051633250574960e11b81526001600160a01b0385166004820152602401610431565b505b5050505050565b6060610e3082610a9f565b505f610e4660408051602081019091525f815290565b90505f815111610e645760405180602001604052805f815250610e8f565b80610e6e84611071565b604051602001610e7f929190611784565b6040516020818303038152906040525b9392505050565b610408828260405180602001604052805f815250611100565b5f828152600660205260409020610ec68282611857565b506040518281527ff8e1a15aba9398e019f0b49df1a4fde98ee17ae345cb5f6b5e2c27f5033e8ce79060200160405180910390a15050565b610651848484611117565b8080610f1d57506001600160a01b03821615155b15610fde575f610f2c84610a9f565b90506001600160a01b03831615801590610f585750826001600160a01b0316816001600160a01b031614155b8015610f6b5750610f69818461099f565b155b15610f945760405163a9fbf51f60e01b81526001600160a01b0384166004820152602401610431565b8115610fdc5783856001600160a01b0316826001600160a01b03167f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b92560405160405180910390a45b505b50505f90815260046020526040902080546001600160a01b0319166001600160a01b0392909216919091179055565b6110188383836111c4565b610342576001600160a01b03831661104657604051637e27328960e01b815260048101829052602401610431565b60405163177e802f60e01b81526001600160a01b038316600482015260248101829052604401610431565b60605f61107d83611225565b60010190505f816001600160401b0381111561109b5761109b6114dd565b6040519080825280601f01601f1916602001820160405280156110c5576020820181803683370190505b5090508181016020015b5f19016f181899199a1a9b1b9c1cb0b131b232b360811b600a86061a8153600a85049450846110cf57509392505050565b61110a83836112fc565b610342335f858585610cfd565b6001600160a01b03821661114057604051633250574960e11b81525f6004820152602401610431565b5f61114c83835f610ae4565b90506001600160a01b03811661117857604051637e27328960e01b815260048101839052602401610431565b836001600160a01b0316816001600160a01b031614610494576040516364283d7b60e01b81526001600160a01b0380861660048301526024820184905282166044820152606401610431565b5f6001600160a01b038316158015906108115750826001600160a01b0316846001600160a01b031614806111fd57506111fd848461099f565b806108115750505f908152600460205260409020546001600160a01b03908116911614919050565b5f8072184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b83106112635772184f03e93ff9f4daa797ed6e38ed64bf6a1f0160401b830492506040015b6d04ee2d6d415b85acef8100000000831061128f576d04ee2d6d415b85acef8100000000830492506020015b662386f26fc1000083106112ad57662386f26fc10000830492506010015b6305f5e10083106112c5576305f5e100830492506008015b61271083106112d957612710830492506004015b606483106112eb576064830492506002015b600a831061031d5760010192915050565b6001600160a01b03821661132557604051633250574960e11b81525f6004820152602401610431565b5f61133183835f610ae4565b90506001600160a01b03811615610342576040516339e3563760e11b81525f6004820152602401610431565b6001600160e01b031981168114610a06575f5ffd5b5f60208284031215611382575f5ffd5b8135610e8f8161135d565b80356001600160a01b03811681146113a3575f5ffd5b919050565b5f5f604083850312156113b9575f5ffd5b823591506113c96020840161138d565b90509250929050565b5f81518084528060208401602086015e5f602082860101526020601f19601f83011685010191505092915050565b602081525f610e8f60208301846113d2565b5f60208284031215611422575f5ffd5b5035919050565b5f5f6040838503121561143a575f5ffd5b6114438361138d565b946020939093013593505050565b5f5f5f60608486031215611463575f5ffd5b61146c8461138d565b925061147a6020850161138d565b929592945050506040919091013590565b5f6020828403121561149b575f5ffd5b610e8f8261138d565b5f5f604083850312156114b5575f5ffd5b6114be8361138d565b9150602083013580151581146114d2575f5ffd5b809150509250929050565b634e487b7160e01b5f52604160045260245ffd5b5f5f5f5f60808587031215611504575f5ffd5b61150d8561138d565b935061151b6020860161138d565b92506040850135915060608501356001600160401b0381111561153c575f5ffd5b8501601f8101871361154c575f5ffd5b80356001600160401b03811115611565576115656114dd565b604051601f8201601f19908116603f011681016001600160401b0381118282101715611593576115936114dd565b6040528181528282016020018910156115aa575f5ffd5b816020840160208301375f6020838301015280935050505092959194509250565b634e487b7160e01b5f52602160045260245ffd5b600381106115fb57634e487b7160e01b5f52602160045260245ffd5b9052565b5f6080820190508251825261ffff60208401511660208301526001600160401b036040840151166040830152606083015161163d60608401826115df565b5092915050565b5f5f5f5f5f5f60a08789031215611659575f5ffd5b6116628761138d565b955060208701359450604087013561ffff8116811461167f575f5ffd5b935060608701356001600160401b038116811461169a575f5ffd5b925060808701356001600160401b038111156116b4575f5ffd5b8701601f810189136116c4575f5ffd5b80356001600160401b038111156116d9575f5ffd5b8960208284010111156116ea575f5ffd5b60208201935080925050509295509295509295565b5f5f60408385031215611710575f5ffd5b6117198361138d565b91506113c96020840161138d565b600181811c9082168061173b57607f821691505b60208210810361175957634e487b7160e01b5f52602260045260245ffd5b50919050565b6020810161031d82846115df565b5f81518060208401855e5f93019283525090919050565b5f610811611792838661176d565b8461176d565b5f600182016117b557634e487b7160e01b5f52601160045260245ffd5b5060010190565b6001600160a01b03858116825284166020820152604081018390526080606082018190525f906117ee908301846113d2565b9695505050505050565b5f60208284031215611808575f5ffd5b8151610e8f8161135d565b601f82111561034257805f5260205f20601f840160051c810160208510156118385750805b601f840160051c820191505b81811015610e1e575f8155600101611844565b81516001600160401b03811115611870576118706114dd565b6118848161187e8454611727565b84611813565b6020601f8211600181146118b6575f831561189f5750848201515b5f19600385901b1c1916600184901b178455610e1e565b5f84815260208120601f198516915b828110156118e557878501518255602094850194600190920191016118c5565b508482101561190257868401515f19600387901b60f8161c191681555b50505050600190811b0190555056fea2646970667358221220c022143627a08a86da67414f0c9eae28be12c929b491bd1e378deb2f5bee5ae464736f6c634300081c0033",
}

// LoanNoteABI is the input ABI used to generate the binding from.
// Deprecated: Use LoanNoteMetaData.ABI instead.
var LoanNoteABI = LoanNoteMetaData.ABI

// LoanNoteBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use LoanNoteMetaData.Bin instead.
var LoanNoteBin = LoanNoteMetaData.Bin

// DeployLoanNote deploys a new Ethereum contract, binding an instance of LoanNote to it.
func DeployLoanNote(auth *bind.TransactOpts, backend bind.ContractBackend, initialOwner common.Address) (common.Address, *types.Transaction, *LoanNote, error) {
	parsed, err := LoanNoteMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(LoanNoteBin), backend, initialOwner)
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

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LoanNote *LoanNoteCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _LoanNote.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LoanNote *LoanNoteSession) Owner() (common.Address, error) {
	return _LoanNote.Contract.Owner(&_LoanNote.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_LoanNote *LoanNoteCallerSession) Owner() (common.Address, error) {
	return _LoanNote.Contract.Owner(&_LoanNote.CallOpts)
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

// AdminTransfer is a paid mutator transaction binding the contract method 0x0483fc5e.
//
// Solidity: function adminTransfer(uint256 tokenId, address to) returns()
func (_LoanNote *LoanNoteTransactor) AdminTransfer(opts *bind.TransactOpts, tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "adminTransfer", tokenId, to)
}

// AdminTransfer is a paid mutator transaction binding the contract method 0x0483fc5e.
//
// Solidity: function adminTransfer(uint256 tokenId, address to) returns()
func (_LoanNote *LoanNoteSession) AdminTransfer(tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _LoanNote.Contract.AdminTransfer(&_LoanNote.TransactOpts, tokenId, to)
}

// AdminTransfer is a paid mutator transaction binding the contract method 0x0483fc5e.
//
// Solidity: function adminTransfer(uint256 tokenId, address to) returns()
func (_LoanNote *LoanNoteTransactorSession) AdminTransfer(tokenId *big.Int, to common.Address) (*types.Transaction, error) {
	return _LoanNote.Contract.AdminTransfer(&_LoanNote.TransactOpts, tokenId, to)
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

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LoanNote *LoanNoteTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LoanNote *LoanNoteSession) RenounceOwnership() (*types.Transaction, error) {
	return _LoanNote.Contract.RenounceOwnership(&_LoanNote.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_LoanNote *LoanNoteTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _LoanNote.Contract.RenounceOwnership(&_LoanNote.TransactOpts)
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

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LoanNote *LoanNoteTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _LoanNote.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LoanNote *LoanNoteSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _LoanNote.Contract.TransferOwnership(&_LoanNote.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_LoanNote *LoanNoteTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _LoanNote.Contract.TransferOwnership(&_LoanNote.TransactOpts, newOwner)
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

// LoanNoteOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the LoanNote contract.
type LoanNoteOwnershipTransferredIterator struct {
	Event *LoanNoteOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *LoanNoteOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(LoanNoteOwnershipTransferred)
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
		it.Event = new(LoanNoteOwnershipTransferred)
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
func (it *LoanNoteOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *LoanNoteOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// LoanNoteOwnershipTransferred represents a OwnershipTransferred event raised by the LoanNote contract.
type LoanNoteOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LoanNote *LoanNoteFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*LoanNoteOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LoanNote.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &LoanNoteOwnershipTransferredIterator{contract: _LoanNote.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LoanNote *LoanNoteFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *LoanNoteOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _LoanNote.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(LoanNoteOwnershipTransferred)
				if err := _LoanNote.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_LoanNote *LoanNoteFilterer) ParseOwnershipTransferred(log types.Log) (*LoanNoteOwnershipTransferred, error) {
	event := new(LoanNoteOwnershipTransferred)
	if err := _LoanNote.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
