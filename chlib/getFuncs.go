package chlib

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

func LoadGenericJsonFromFile(path string) (b GenericJson, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	err = json.NewDecoder(file).Decode(&b)
	return
}

func GetCmdRequestJson(client *Client, kind, name string) (ret GenericJson, err error) {
	switch kind {
	case "ns", "namespaces", "namespace":
		apiResult, err := client.GetNameSpaces(name)
		if err != nil {
			return ret, err
		}
		ret = GenericJson(apiResult)
	}
	return
}

type PrettyPrintConfig struct {
	Columns []string
	Data    [][]string
}

type NsResult struct {
	Data struct {
		Metadata struct {
			CreatedAt time.Time `json:"creationTimestamp"`
			Namespace string    `json:"namespace,omitempty"`
		} `json:"metadata"`
		Status struct {
			Hard struct {
				LimitsCpu string `json:"limits.cpu"`
				LimitsMem string `json:"limits.memory"`
			} `json:"hard"`
			Used struct {
				LimitsCpu string `json:"limits.cpu"`
				LimitsMem string `json:"limits.memory"`
			} `json:"used"`
		} `json:"status"`
	} `json:"data"`
}

type nsResponse struct {
	Results []NsResult `json:"results"`
}

func ExtractNsResults(data GenericJson) (res []NsResult, err error) {
	b, _ := json.Marshal(data)
	var resp nsResponse
	if err := json.Unmarshal(b, &resp); err != nil {
		return res, fmt.Errorf("invalid namespace response: %s", err)
	}
	for _, v := range resp.Results {
		if v.Data.Metadata.Namespace != "" {
			res = append(res, v)
		}
	}
	return res, nil
}

func FormatNamespacePrettyPrint(data []NsResult) (ppc PrettyPrintConfig, err error) {
	ppc.Columns = []string{"NAME", "HARD CPU", "HARD MEMORY", "USED CPU", "USED MEMORY", "AGE"}
	for _, v := range data {
		row := []string{
			v.Data.Metadata.Namespace,
			v.Data.Status.Hard.LimitsCpu,
			v.Data.Status.Hard.LimitsMem,
			v.Data.Status.Used.LimitsCpu,
			v.Data.Status.Used.LimitsMem,
			fmt.Sprintf("%dd", int(time.Now().Sub(v.Data.Metadata.CreatedAt).Hours()/24)),
		}
		ppc.Data = append(ppc.Data, row)
	}
	return
}

func PrettyPrint(ppc PrettyPrintConfig, writer io.Writer) {
	table := tablewriter.NewWriter(writer)
	table.SetHeader(ppc.Columns)
	table.AppendBulk(ppc.Data)
	table.Render()
}
