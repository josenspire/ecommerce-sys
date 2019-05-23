package commons

import (
	"crypto/ecdsa"
	. "ecommerce-sys/utils"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"log"
	"net/http"
)

type RequestModel struct {
	Data      string `json:"data"`
	SecretKey string `json:"secretKey"`
	Signature string `json:"signature"`
}

type IAspectControl interface {
	HandleRequest(ct *context.Context)
	HandleResponse(ct *context.Context)
	HandleRequestWithoutEcdh(ct *context.Context)
	HandleResponseWithoutEcdh(ct *context.Context)
}

type AspectControl struct {
	IAspectControl
}
type AspectControlWithoutEcdh struct {
	IAspectControl
}

func (asp *AspectControl) HandleRequest(ct *context.Context) {
	var requestModel *RequestModel

	inputContent := ct.Input.RequestBody
	if len(inputContent) != 0 {
		err := json.Unmarshal(ct.Input.RequestBody, requestModel)
		if err != nil {
			beego.Error(err.Error())
			asp.HandleResponse(ct)
		}

		reqBytes, err := asp.HandleRequestBody(requestModel)
		if err == ErrSignatureInvalid {
			ct.Abort(http.StatusForbidden, err.Error())
		} else if err != nil {
			// TODO: should build completed response body then return to client
			ct.Abort(http.StatusOK, err.Error())
		}

		ct.Input.RequestBody = reqBytes
		log.Println("[Request With ECDH] request body: ", ct.Input.RequestBody)
	}
}

func (asp *AspectControl) HandleResponse(ct *context.Context) {
	// 	TODO: response
	// resArgs := make(map[string]interface{})
	// log.Println("response body: ", ct.Output.JSON(&resArgs, true, true))
}

func (asp *AspectControl) HandleRequestBody(requestModel *RequestModel) (reqBytes []byte, err error) {
	var pubKey *EllipticPublicKey
	var ellipticECDH = &EllipticECDH{}
	var ecdsaPublicKey *ecdsa.PublicKey
	var signatureData *SignatureData
	var verifyResult bool

	ellipticECDH, err = ellipticECDH.ParseECPrivateKeyFromPEM("./../pem/ecdh_priv.pem")
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	// verify signature
	signatureData, err = HandleSignatureData(requestModel.Data, requestModel.Signature)
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	verifyResult = ellipticECDH.VerifySignature(signatureData, ecdsaPublicKey)
	if !verifyResult {
		return nil, ErrSignatureInvalid
	}
	// secret compute and decryption
	pubKey, ecdsaPublicKey, err = ellipticECDH.ParseECPublicKeyFromPEM(requestModel.SecretKey)
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	secret, err := ellipticECDH.ComputeSecret(ellipticECDH.PrivateKey, pubKey)
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	requestData, err := AESDecrypt(requestModel.Data, base64.StdEncoding.EncodeToString(secret))
	if err != nil {
		beego.Error(err.Error())
		return nil, err
	}
	return []byte(requestData), nil
}

func (asp *AspectControlWithoutEcdh) HandleRequestWithoutEcdh(ct *context.Context) {
	log.Println("[Normal Request] request body: ", string(ct.Input.RequestBody))
}

func (asp *AspectControlWithoutEcdh) HandleResponseWithoutEcdh(ct *context.Context) {
	fmt.Print(ct.Output)
	log.Println("[Normal Response] response body: ")
}
