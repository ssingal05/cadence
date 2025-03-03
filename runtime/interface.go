/*
 * Cadence - The resource-oriented smart contract programming language
 *
 * Copyright 2019-2020 Dapper Labs, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package runtime

import (
	"time"

	opentracing "github.com/opentracing/opentracing-go"

	"github.com/onflow/atree"

	"github.com/onflow/cadence"
	"github.com/onflow/cadence/runtime/common"
	"github.com/onflow/cadence/runtime/interpreter"
)

type Interface interface {
	// ResolveLocation resolves an import location.
	ResolveLocation(identifiers []Identifier, location Location) ([]ResolvedLocation, error)
	// GetCode returns the code at a given location
	GetCode(location Location) ([]byte, error)
	// GetProgram attempts gets the program for the given location, if available.
	//
	// NOTE: During execution, this function must always return the *same* program,
	// i.e. it may NOT return a different program,
	// an elaboration in the program that is not annotating the AST in the program;
	// or a program/elaboration and then nothing in a subsequent call.
	//
	// This function must also return what was set using SetProgram,
	// it may NOT return something different or nothing (!) after SetProgram was called.
	//
	// This is not a caching function!
	//
	GetProgram(Location) (*interpreter.Program, error)
	// SetProgram sets the program for the given location.
	SetProgram(Location, *interpreter.Program) error
	// GetValue gets a value for the given key in the storage, owned by the given account.
	GetValue(owner, key []byte) (value []byte, err error)
	// SetValue sets a value for the given key in the storage, owned by the given account.
	SetValue(owner, key, value []byte) (err error)
	// ValueExists returns true if the given key exists in the storage, owned by the given account.
	ValueExists(owner, key []byte) (exists bool, err error)
	// AllocateStorageIndex allocates a new storage index under the given account.
	AllocateStorageIndex(owner []byte) (atree.StorageIndex, error)
	// CreateAccount creates a new account.
	CreateAccount(payer Address) (address Address, err error)
	// AddEncodedAccountKey appends an encoded key to an account.
	AddEncodedAccountKey(address Address, publicKey []byte) error
	// RevokeEncodedAccountKey removes a key from an account by index, add returns the encoded key.
	RevokeEncodedAccountKey(address Address, index int) (publicKey []byte, err error)
	// AddAccountKey appends a key to an account.
	AddAccountKey(address Address, publicKey *PublicKey, hashAlgo HashAlgorithm, weight int) (*AccountKey, error)
	// GetAccountKey retrieves a key from an account by index.
	GetAccountKey(address Address, index int) (*AccountKey, error)
	// RevokeAccountKey removes a key from an account by index.
	RevokeAccountKey(address Address, index int) (*AccountKey, error)
	// UpdateAccountContractCode updates the code associated with an account contract.
	UpdateAccountContractCode(address Address, name string, code []byte) (err error)
	// GetAccountContractCode returns the code associated with an account contract.
	GetAccountContractCode(address Address, name string) (code []byte, err error)
	// RemoveAccountContractCode removes the code associated with an account contract.
	RemoveAccountContractCode(address Address, name string) (err error)
	// GetSigningAccounts returns the signing accounts.
	GetSigningAccounts() ([]Address, error)
	// ProgramLog logs program logs.
	ProgramLog(string) error
	// EmitEvent is called when an event is emitted by the runtime.
	EmitEvent(cadence.Event) error
	// GenerateUUID is called to generate a UUID.
	GenerateUUID() (uint64, error)
	// GetComputationLimit returns the computation limit. A value <= 0 means there is no limit
	GetComputationLimit() uint64
	// SetComputationUsed reports the amount of computation used.
	SetComputationUsed(used uint64) error
	// DecodeArgument decodes a transaction argument against the given type.
	DecodeArgument(argument []byte, argumentType cadence.Type) (cadence.Value, error)
	// GetCurrentBlockHeight returns the current block height.
	GetCurrentBlockHeight() (uint64, error)
	// GetBlockAtHeight returns the block at the given height.
	GetBlockAtHeight(height uint64) (block Block, exists bool, err error)
	// UnsafeRandom returns a random uint64, where the process of random number derivation is not cryptographically
	// secure.
	UnsafeRandom() (uint64, error)
	// VerifySignature returns true if the given signature was produced by signing the given tag + data
	// using the given public key, signature algorithm, and hash algorithm.
	VerifySignature(
		signature []byte,
		tag string,
		signedData []byte,
		publicKey []byte,
		signatureAlgorithm SignatureAlgorithm,
		hashAlgorithm HashAlgorithm,
	) (bool, error)
	// Hash returns the digest of hashing the given data with using the given hash algorithm
	Hash(data []byte, tag string, hashAlgorithm HashAlgorithm) ([]byte, error)
	// GetAccountBalance gets accounts default flow token balance.
	GetAccountBalance(address common.Address) (value uint64, err error)
	// GetAccountAvailableBalance gets accounts default flow token balance - balance that is reserved for storage.
	GetAccountAvailableBalance(address common.Address) (value uint64, err error)
	// GetStorageUsed gets storage used in bytes by the address at the moment of the function call.
	GetStorageUsed(address Address) (value uint64, err error)
	// GetStorageCapacity gets storage capacity in bytes on the address.
	GetStorageCapacity(address Address) (value uint64, err error)
	// ImplementationDebugLog logs implementation log statements on a debug-level
	ImplementationDebugLog(message string) error
	// ValidatePublicKey verifies the validity of a public key.
	ValidatePublicKey(key *PublicKey) (bool, error)
	// GetAccountContractNames returns the names of all contracts deployed in an account.
	GetAccountContractNames(address Address) ([]string, error)
	// RecordTrace records a opentracing trace
	RecordTrace(operation string, location common.Location, duration time.Duration, logs []opentracing.LogRecord)
	// BLSVerifyPOP verifies a proof of possession (PoP) for the receiver public key.
	BLSVerifyPOP(pk *PublicKey, s []byte) (bool, error)
	// AggregateBLSSignatures aggregates multiple BLS signatures into one.
	AggregateBLSSignatures(sigs [][]byte) ([]byte, error)
	// AggregateBLSPublicKeys aggregates multiple BLS public keys into one.
	AggregateBLSPublicKeys(keys []*PublicKey) (*PublicKey, error)
	// ResourceOwnerChanged gets called when a resource's owner changed (if enabled)
	ResourceOwnerChanged(resource *interpreter.CompositeValue, oldOwner common.Address, newOwner common.Address)
}

type Metrics interface {
	ProgramParsed(location common.Location, duration time.Duration)
	ProgramChecked(location common.Location, duration time.Duration)
	ProgramInterpreted(location common.Location, duration time.Duration)
}
