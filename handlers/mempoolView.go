package handlers

import (
	"eth2-exporter/services"
	"eth2-exporter/templates"
	"eth2-exporter/types"
	"eth2-exporter/utils"
	"fmt"
	"html/template"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
)

func MempoolView(w http.ResponseWriter, r *http.Request) {
	mempool := services.LatestMempoolTransactions()
	formatedData := formatToTable(mempool)

	var err error
	var mempoolViewTemplate = templates.GetTemplate("layout.html", "mempoolview.html")

	w.Header().Set("Content-Type", "text/html")
	data := InitPageData(w, r, "services", "/mempool", "Pending Mempool Transactions")

	data.Data = formatedData

	err = mempoolViewTemplate.ExecuteTemplate(w, "layout", data)
	if err != nil {
		logger.Errorf("error executing template for %v route: %v", r.URL.String(), err)
		http.Error(w, "Internal server error", 503)
		return
	}

}

// This is a helper function. It replaces Nil or empty receiver Address with a string in case case of a new contract creation.
// This function catches the Nil exception
func _isContractCreation(tx *common.Address) string {
	if tx == nil {
		return "Contract Creation"
	}
	return string(utils.FormatAddressAll(tx.Bytes(), "", false, "address", "", int(12), int(12), true))
}

// This Function formats each Transaction into Html string.
// This makes all calculations faster, reducing browser's rendering time.
func formatToTable(content *types.RawMempoolResponse) *types.DataTableResponse {
	dataTable := &types.DataTableResponse{}

	for _, pendingData := range content.Pending {
		for _, tx := range pendingData {
			dataTable.Data = append(dataTable.Data, []any{
				// TODO: link to tx page once it has a mempool view implemented
				template.HTML(fmt.Sprintf(`<span class="text-monospace" >%v</span>`, tx.Hash.String())),
				utils.FormatAddressAll(tx.From.Bytes(), "", false, "address", "", int(12), int(12), true),
				_isContractCreation(tx.To),
				utils.FormatAmount((*big.Int)(tx.Value), "ETH", 5),
				utils.FormatAddCommasFormated(float64(tx.Gas.ToInt().Int64()), 0),
				utils.FormatAmountFormated(tx.GasPrice.ToInt(), "GWei", 5, 0, true, true, false),
				tx.Nonce.ToInt(),
			})
		}
	}
	return dataTable
}

// type formatedTx struct {
// 	Hash      template.HTML `json:"hash"`
// 	From      template.HTML `json:"from"`
// 	To        template.HTML `default:"Empty address"`
// 	Value     template.HTML `json:"value"`
// 	Gas       template.HTML `json:"gas"`
// 	GasFeeCap template.HTML `json:"maxFeePerGas,omitempty"`
// 	GasPrice  template.HTML `json:"gasPrice"`
// 	Nonce     template.HTML `json:"nonce"`
// }
