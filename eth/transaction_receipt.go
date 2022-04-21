package eth

import (
	"errors"
	"fmt"
	"strings"

	"github.com/INFURA/go-ethlibs/rlp"
)

type TransactionReceipt struct {
	Type              *Quantity `json:"type,omitempty"`
	TransactionHash   Hash      `json:"transactionHash"`
	TransactionIndex  Quantity  `json:"transactionIndex"`
	BlockHash         Hash      `json:"blockHash"`
	BlockNumber       Quantity  `json:"blockNumber"`
	From              Address   `json:"from"`
	To                *Address  `json:"to"`
	CumulativeGasUsed Quantity  `json:"cumulativeGasUsed"`
	GasUsed           Quantity  `json:"gasUsed"`
	ContractAddress   *Address  `json:"contractAddress"`
	Logs              []Log     `json:"logs"`
	LogsBloom         Data256   `json:"logsBloom"`
	Root              *Data32   `json:"root,omitempty"`
	Status            *Quantity `json:"status,omitempty"`
	EffectiveGasPrice *Quantity `json:"effectiveGasPrice,omitempty"`
}

// TransactionType returns the transactions EIP-2718 type, or TransactionTypeLegacy for pre-EIP-2718 transactions.
func (t *TransactionReceipt) TransactionType() int64 {
	if t.Type == nil {
		return TransactionTypeLegacy
	}

	return t.Type.Int64()
}

// RequiredFields inspects the Transaction Type and returns an error if any required fields are missing
func (t *TransactionReceipt) RequiredFields() error {
	var fields []string
	switch t.TransactionType() {
	case TransactionTypeLegacy:
		// LegacyReceipt is rlp([status, cumulativeGasUsed, logsBloom, logs])
		// only .Status is a pointer at the moment
		if t.Status == nil {
			fields = append(fields, "status")
		}
		return nil
	case TransactionTypeAccessList:
		// The EIP-2718 ReceiptPayload for this transaction is rlp([status, cumulativeGasUsed, logsBloom, logs]).
		// Same as TransactionTypeLegacy.
		if t.Status == nil {
			fields = append(fields, "status")
		}
		return nil
	case TransactionTypeDynamicFee:
		// The EIP-2718 ReceiptPayload for this transaction is rlp([status, cumulative_transaction_gas_used, logs_bloom, logs]).
		// Same as TransactionTypeLegacy.
		if t.Status == nil {
			fields = append(fields, "status")
		}
	}

	if len(fields) > 0 {
		return fmt.Errorf("missing required field(s) %s for transaction type", strings.Join(fields, ","))
	}

	return nil
}

func (t *TransactionReceipt) RawRepresentation() (*Data, error) {
	if err := t.RequiredFields(); err != nil {
		return nil, err
	}

	logsRLP := func() rlp.Value {
		list := make([]rlp.Value, len(t.Logs))
		for i := range t.Logs {
			list[i] = t.Logs[i].RLP()
		}
		return rlp.Value{List: list}
	}

	switch t.TransactionType() {
	case TransactionTypeLegacy:
		// LegacyReceipt is rlp([status, cumulativeGasUsed, logsBloom, logs])
		message := rlp.Value{List: []rlp.Value{
			t.Status.RLP(),
			t.CumulativeGasUsed.RLP(),
			t.LogsBloom.RLP(),
			logsRLP(),
		}}
		if encoded, err := message.Encode(); err != nil {
			return nil, err
		} else {
			return NewData(encoded)
		}
	case TransactionTypeAccessList:
		// The EIP-2718 ReceiptPayload for this transaction is rlp([status, cumulativeGasUsed, logsBloom, logs]).
		// Same as TransactionTypeLegacy.
		typePrefix, err := t.Type.RLP().Encode()
		if err != nil {
			return nil, err
		}
		payload := rlp.Value{List: []rlp.Value{
			t.Status.RLP(),
			t.CumulativeGasUsed.RLP(),
			t.LogsBloom.RLP(),
			logsRLP(),
		}}
		if encodedPayload, err := payload.Encode(); err != nil {
			return nil, err
		} else {
			return NewData(typePrefix + encodedPayload[2:])
		}
	case TransactionTypeDynamicFee:
		// The EIP-2718 ReceiptPayload for this transaction is rlp([status, cumulative_transaction_gas_used, logs_bloom, logs]).
		// Same as TransactionTypeLegacy.
		typePrefix, err := t.Type.RLP().Encode()
		if err != nil {
			return nil, err
		}
		payload := rlp.Value{List: []rlp.Value{
			t.Status.RLP(),
			t.CumulativeGasUsed.RLP(),
			t.LogsBloom.RLP(),
			logsRLP(),
		}}
		if encodedPayload, err := payload.Encode(); err != nil {
			return nil, err
		} else {
			return NewData(typePrefix + encodedPayload[2:])
		}
	default:
		return nil, errors.New("unsupported transaction type")
	}
}
