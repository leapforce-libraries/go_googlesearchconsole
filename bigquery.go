package googlesearchconsole

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/api/iterator"
)

const tableRefreshToken string = "leapforce.refreshtokens"
const api string = "GoogleSearchConsole"

// BigQueryGetRefreshToken get refreshtoken from BigQuery
//
func (gsc *GoogleSearchConsole) GetTokenFromBigQuery() error {
	fmt.Println("***GetTokenFromBigQuery***")
	// create client
	bqClient, err := gsc.BigQuery.CreateClient()
	if err != nil {
		fmt.Println("\nerror in BigQueryCreateClient")
		return err
	}

	ctx := context.Background()

	//sql := "SELECT refreshtoken AS RefreshToken FROM `" + tableRefreshToken + "` WHERE client_id = '" + gsc.ClientID + "'"
	sql := fmt.Sprintf("SELECT refreshtoken AS RefreshToken FROM `%s` WHERE api = '%s' AND client_id = '%s'", tableRefreshToken, api, gsc.ClientID)

	//fmt.Println(sql)

	q := bqClient.Query(sql)
	it, err := q.Read(ctx)
	if err != nil {
		return err
	}

	token := new(Token)

	for {
		err := it.Next(token)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}

		break
	}

	if gsc.Token == nil {
		gsc.Token = new(Token)
	}

	gsc.Token.TokenType = "bearer"
	gsc.Token.Expiry = time.Now().Add(-10 * time.Second)
	gsc.Token.RefreshToken = token.RefreshToken
	gsc.Token.AccessToken = ""

	return nil
}

// BigQuerySaveToken saves refreshtoken to BigQuery
//
func (gsc *GoogleSearchConsole) SaveTokenToBigQuery() error {
	// create client
	bqClient, err := gsc.BigQuery.CreateClient()
	if err != nil {
		fmt.Println("\nerror in BigQueryCreateClient")
		return err
	}

	ctx := context.Background()

	sql := "MERGE `" + tableRefreshToken + "` AS TARGET " +
		"USING  (SELECT '" + api + "' AS api,'" + gsc.ClientID + "' AS client_id,'" + gsc.Token.RefreshToken + "' AS refreshtoken) AS SOURCE " +
		" ON TARGET.api = SOURCE.api " +
		" AND TARGET.client_id = SOURCE.client_id " +
		"WHEN MATCHED THEN " +
		"	UPDATE " +
		"	SET refreshtoken = SOURCE.refreshtoken " +
		"WHEN NOT MATCHED BY TARGET THEN " +
		"	INSERT (api, client_id, refreshtoken) " +
		"	VALUES (SOURCE.api, SOURCE.client_id, SOURCE.refreshtoken)"

	q := bqClient.Query(sql)

	job, err := q.Run(ctx)
	if err != nil {
		return err
	}

	for {
		status, err := job.Status(ctx)
		if err != nil {
			return err
		}
		if status.Done() {
			if status.Err() != nil {
				return status.Err()
			}
			break
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}
