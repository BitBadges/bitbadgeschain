package lifecycle

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/bitbadges/bitbadgeschain/app"
	"github.com/bitbadges/bitbadgeschain/app/params"
	storetypes "github.com/cosmos/cosmos-sdk/store/v2/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const ProtocolVersion = 1

var runMu sync.Mutex // The app's EVM configuration is process-global.
var ChainCommit = "unknown"

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, s := range info.Settings {
			if s.Key == "vcs.revision" {
				ChainCommit = s.Value
			}
		}
	}
}

type Actor struct {
	Name    string    `json:"name"`
	Address string    `json:"address"`
	Coins   sdk.Coins `json:"coins"`
}
type Message struct {
	TypeURL string          `json:"typeUrl"`
	Value   json.RawMessage `json:"value"`
}
type Expect struct {
	Success       bool   `json:"success"`
	ErrorContains string `json:"errorContains,omitempty"`
}
type Step struct {
	ID            string      `json:"id"`
	Actor         string      `json:"actor,omitempty"`
	Message       *Message    `json:"message,omitempty"`
	Expect        *Expect     `json:"expect,omitempty"`
	AdvanceTimeMs string      `json:"advanceTimeMs,omitempty"`
	Assertions    []Assertion `json:"assertions"`
}
type Scenario struct {
	Version int     `json:"version"`
	ID      string  `json:"id"`
	TimeMs  string  `json:"timeMs"`
	Actors  []Actor `json:"actors"`
	Steps   []Step  `json:"steps"`
}
type Assertion struct {
	Kind               string          `json:"kind"`
	CollectionID       string          `json:"collectionId,omitempty"`
	Address            string          `json:"address,omitempty"`
	Denom              string          `json:"denom,omitempty"`
	Expected           json.RawMessage `json:"expected"`
	TokenID            string          `json:"tokenId,omitempty"`
	OwnershipTime      string          `json:"ownershipTime,omitempty"`
	ApprovalID         string          `json:"approvalId,omitempty"`
	AmountTrackerID    string          `json:"amountTrackerId,omitempty"`
	Level              string          `json:"level,omitempty"`
	TrackerType        string          `json:"trackerType,omitempty"`
	AddressForApproval string          `json:"addressForApproval,omitempty"`
}
type AssertionResult struct {
	Kind     string          `json:"kind"`
	Passed   bool            `json:"passed"`
	Expected json.RawMessage `json:"expected"`
	Actual   json.RawMessage `json:"actual"`
}
type StepResult struct {
	ID         string            `json:"id"`
	Passed     bool              `json:"passed"`
	Success    bool              `json:"success"`
	Error      string            `json:"error,omitempty"`
	TimeMs     string            `json:"timeMs"`
	Events     sdk.Events        `json:"events"`
	Assertions []AssertionResult `json:"assertions"`
}
type Coverage struct {
	Excluded []string `json:"excluded"`
}
type Result struct {
	Version      int          `json:"version"`
	ID           string       `json:"id"`
	Passed       bool         `json:"passed"`
	Execution    string       `json:"execution"`
	ChainCommit  string       `json:"chainCommit"`
	ScenarioHash string       `json:"scenarioHash"`
	Steps        []StepResult `json:"steps"`
	Coverage     Coverage     `json:"coverage"`
}

func strictJSON(data []byte, target any) error {
	d := json.NewDecoder(bytes.NewReader(data))
	d.DisallowUnknownFields()
	if err := d.Decode(target); err != nil {
		return err
	}
	var extra any
	if err := d.Decode(&extra); err != io.EOF {
		return fmt.Errorf("expected exactly one JSON value")
	}
	return nil
}
func positiveMs(s string) (int64, error) {
	v, e := strconv.ParseInt(s, 10, 64)
	if e != nil || v < 1 {
		return 0, fmt.Errorf("invalid positive milliseconds %q", s)
	}
	return v, nil
}

func RunJSON(data []byte) (out *Result, runErr error) {
	defer func() {
		if value := recover(); value != nil {
			out = nil
			runErr = fmt.Errorf("lifecycle execution panic: %v", value)
		}
	}()
	if len(data) > 8<<20 {
		return nil, fmt.Errorf("scenario exceeds 8 MiB")
	}
	var s Scenario
	if err := strictJSON(data, &s); err != nil {
		return nil, err
	}
	now, err := positiveMs(s.TimeMs)
	if err != nil {
		return nil, err
	}
	if s.Version != ProtocolVersion || s.ID == "" || len(s.Steps) == 0 || len(s.Steps) > 1000 || len(s.Actors) == 0 || len(s.Actors) > 100 {
		return nil, fmt.Errorf("invalid scenario version, id, actors or steps")
	}
	runMu.Lock()
	defer runMu.Unlock()
	params.InitSDKConfigWithoutSeal()
	actors := map[string]string{}
	addresses := map[string]bool{}
	for _, a := range s.Actors {
		if a.Name == "" || actors[a.Name] != "" || addresses[a.Address] {
			return nil, fmt.Errorf("duplicate/empty actor")
		}
		if _, err := sdk.AccAddressFromBech32(a.Address); err != nil {
			return nil, err
		}
		if err := a.Coins.Validate(); err != nil {
			return nil, err
		}
		actors[a.Name] = a.Address
		addresses[a.Address] = true
	}
	ids := map[string]bool{}
	for _, step := range s.Steps {
		if step.ID == "" || ids[step.ID] || len(step.Assertions) == 0 {
			return nil, fmt.Errorf("step needs unique id and non-vacuous assertions")
		}
		ids[step.ID] = true
		if (step.Message == nil) == (step.AdvanceTimeMs == "") {
			return nil, fmt.Errorf("step must have exactly one message or time advance")
		}
		if step.Message != nil {
			if step.Expect == nil || actors[step.Actor] == "" {
				return nil, fmt.Errorf("message needs actor and explicit expected outcome")
			}
			if step.Expect.Success && step.Expect.ErrorContains != "" {
				return nil, fmt.Errorf("successful outcome cannot expect error")
			}
		}
		if step.AdvanceTimeMs != "" {
			if _, err := positiveMs(step.AdvanceTimeMs); err != nil {
				return nil, err
			}
			if step.Actor != "" || step.Expect != nil {
				return nil, fmt.Errorf("time step cannot specify actor or outcome")
			}
		}
		for _, a := range step.Assertions {
			if err := validateAssertion(a); err != nil {
				return nil, err
			}
		}
	}
	home, err := os.MkdirTemp("", "bitbadges-lifecycle-")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(home)
	application := app.SetupWithAppOptions(false, map[string]interface{}{"home": home})
	defer application.Close()
	ctx := application.NewContext(false).WithBlockTime(time.UnixMilli(now).UTC()).WithBlockHeight(1)
	for _, a := range s.Actors {
		if len(a.Coins) > 0 {
			addr, _ := sdk.AccAddressFromBech32(a.Address)
			if err := application.BankKeeper.MintCoins(ctx, "mint", a.Coins); err != nil {
				return nil, err
			}
			if err := application.BankKeeper.SendCoinsFromModuleToAccount(ctx, "mint", addr, a.Coins); err != nil {
				return nil, err
			}
		}
	}
	hash := sha256.Sum256(data)
	result := &Result{Version: 1, ID: s.ID, Passed: true, Execution: "module", ChainCommit: ChainCommit, ScenarioHash: hex.EncodeToString(hash[:]), Steps: []StepResult{}, Coverage: Coverage{Excluded: []string{"signatures", "fees", "sequences", "IBC", "external services", "block hooks"}}}
	for _, step := range s.Steps {
		sr := StepResult{ID: step.ID, Passed: true, Success: true, Events: sdk.Events{}, Assertions: []AssertionResult{}}
		before := map[int]sdkmath.Int{}
		for i, a := range step.Assertions {
			if a.Kind == "coinDelta" {
				addr, _ := sdk.AccAddressFromBech32(a.Address)
				before[i] = application.BankKeeper.GetBalance(ctx, addr, a.Denom).Amount
			}
		}
		if step.AdvanceTimeMs != "" {
			delta, _ := positiveMs(step.AdvanceTimeMs)
			if now > 9223372036854775807-delta {
				return nil, fmt.Errorf("time overflow")
			}
			now += delta
			ctx = ctx.WithBlockTime(time.UnixMilli(now).UTC())
		} else {
			msg, err := decodeMessage(application, step.Message)
			if err != nil {
				return nil, fmt.Errorf("step %s: %w", step.ID, err)
			}
			signers, _, err := application.AppCodec().GetMsgV1Signers(msg)
			if err != nil {
				return nil, err
			}
			addr, _ := sdk.AccAddressFromBech32(actors[step.Actor])
			if len(signers) != 1 || !bytes.Equal(signers[0], addr) {
				return nil, fmt.Errorf("step %s actor does not match message signer", step.ID)
			}
			cache, write := ctx.CacheContext()
			cache = cache.WithEventManager(sdk.NewEventManager()).WithGasMeter(storetypes.NewGasMeter(50000000))
			if v, ok := msg.(interface{ ValidateBasic() error }); ok {
				err = v.ValidateBasic()
			}
			if err == nil {
				handler := application.MsgServiceRouter().Handler(msg)
				if handler == nil {
					return nil, fmt.Errorf("unsupported message %s", step.Message.TypeURL)
				}
				var response *sdk.Result
				response, err = handler(cache, msg)
				if err == nil && response != nil {
					for _, event := range response.Events {
						sr.Events = append(sr.Events, sdk.Event(event))
					}
				}
			}
			sr.Success = err == nil
			if err == nil {
				write()
			} else {
				sr.Error = err.Error()
			}
			sr.Passed = sr.Success == step.Expect.Success && (step.Expect.ErrorContains == "" || strings.Contains(sr.Error, step.Expect.ErrorContains))
		}
		sr.TimeMs = strconv.FormatInt(now, 10)
		for i, a := range step.Assertions {
			actual, err := observe(application, ctx, a, before[i])
			if err != nil {
				return nil, fmt.Errorf("step %s assertion: %w", step.ID, err)
			}
			equal, err := compare(a, actual)
			if err != nil {
				return nil, err
			}
			sr.Assertions = append(sr.Assertions, AssertionResult{a.Kind, equal, a.Expected, actual})
			sr.Passed = sr.Passed && equal
		}
		result.Passed = result.Passed && sr.Passed
		result.Steps = append(result.Steps, sr)
	}
	return result, nil
}
func decodeMessage(application *app.App, m *Message) (sdk.Msg, error) {
	if !strings.HasPrefix(m.TypeURL, "/tokenization.") && m.TypeURL != "/cosmos.bank.v1beta1.MsgSend" {
		return nil, fmt.Errorf("unsupported message %s", m.TypeURL)
	}
	var value map[string]json.RawMessage
	if err := json.Unmarshal(m.Value, &value); err != nil || value == nil {
		return nil, fmt.Errorf("message value must be object")
	}
	if _, ok := value["@type"]; ok {
		return nil, fmt.Errorf("value cannot override typeUrl")
	}
	value["@type"], _ = json.Marshal(m.TypeURL)
	data, _ := json.Marshal(value)
	var msg sdk.Msg
	err := application.AppCodec().UnmarshalInterfaceJSON(data, &msg)
	return msg, err
}

func (e *Expect) UnmarshalJSON(data []byte) error {
	var v struct {
		Success       *bool  `json:"success"`
		ErrorContains string `json:"errorContains,omitempty"`
	}
	if err := strictJSON(data, &v); err != nil {
		return err
	}
	if v.Success == nil {
		return fmt.Errorf("expect.success must be explicit")
	}
	e.Success = *v.Success
	e.ErrorContains = v.ErrorContains
	return nil
}
