package tinkoff

import (
	"fmt"
)

const (
	PayTypeOneStep  = "O"
	PayTypeTwoSteps = "T"
)

type InitRequest struct {
	BaseRequest

	Amount          uint64            `json:"Amount,omitempty"`          // Сумма в копейках
	OrderID         string            `json:"OrderId"`                   // Идентификатор заказа в системе продавца
	ClientIP        string            `json:"IP,omitempty"`              // IP-адрес покупателя
	Description     string            `json:"Description,omitempty"`     // Описание заказа
	Language        string            `json:"Language,omitempty"`        // Язык платежной формы: ru или en
	Recurrent       string            `json:"Recurrent,omitempty"`       // Y для регистрации автоплатежа. Можно использовать SetIsRecurrent(true)
	CustomerKey     string            `json:"CustomerKey,omitempty"`     // Идентификатор покупателя в системе продавца. Передается вместе с параметром CardId. См. метод GetCardList
	Data            map[string]string `json:"DATA"`                      // Дополнительные параметры платежа
	Receipt         *Receipt          `json:"Receipt,omitempty"`         // Чек
	RedirectDueDate Time              `json:"RedirectDueDate,omitempty"` // Срок жизни ссылки
	NotificationURL string            `json:"NotificationURL,omitempty"` // Адрес для получения http нотификаций
	SuccessURL      string            `json:"SuccessURL,omitempty"`      // Страница успеха
	FailURL         string            `json:"FailURL,omitempty"`         // Страница ошибки
	PayType         string            `json:"PayType,omitempty"`         // Тип оплаты. см. PayType*
	Shops           *[]Shop           `json:"Shops,omitempty"`           // Объект с данными партнера
}

type Shop struct {
	ShopCode string `json:"ShopCode,omitempty"` // Код магазина. Для параметра ShopСode необходимо использовать значение параметра Submerchant_ID, полученного при регистрации через xml.
	Amount   uint64 `json:"Amount,omitempty"`   // Суммаперечисленияв копейкахпо реквизитам ShopCode за вычетом Fee
	Name     string `json:"Name,omitempty"`     // Наименованиепозиции
	Fee      string `json:"Fee,omitempty"`      // Часть суммы Операции оплаты или % от суммы Операции оплаты. Fee удерживается из возмещения третьего лица (ShopCode) в пользу Предприятия по операциям оплаты.
}

func (i *InitRequest) SetIsRecurrent(r bool) {
	if r {
		i.Recurrent = "Y"
	} else {
		i.Recurrent = ""
	}
}

func (i *InitRequest) GetValuesForToken() map[string]string {
	v := map[string]string{
		"OrderId":         i.OrderID,
		"IP":              i.ClientIP,
		"Description":     i.Description,
		"Language":        i.Language,
		"CustomerKey":     i.CustomerKey,
		"RedirectDueDate": i.RedirectDueDate.String(),
		"NotificationURL": i.NotificationURL,
		"SuccessURL":      i.SuccessURL,
		"FailURL":         i.FailURL,
	}
	serializeUintToMapIfNonEmpty(&v, "Amount", i.Amount)
	return v
}

type InitResponse struct {
	BaseResponse
	Amount     uint64 `json:"Amount"`               // Сумма в копейках
	OrderID    string `json:"OrderId"`              // Номер заказа в системе Продавца
	Status     string `json:"Status"`               // Статус транзакции
	PaymentID  string `json:"PaymentId"`            // Уникальный идентификатор транзакции в системе Банка. По офф. документации это number(20), но фактически значение передается в виде строки.
	PaymentURL string `json:"PaymentURL,omitempty"` // Ссылка на страницу оплаты. По умолчанию ссылка доступна в течении 24 часов.
}

func (c *Client) Init(request *InitRequest) (*InitResponse, error) {
	response, err := c.PostRequest("/Init", request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var res InitResponse
	err = c.decodeResponse(response, &res)
	if err != nil {
		return nil, err
	}

	err = res.Error()
	if res.Status != StatusNew {
		err = errorConcat(err, fmt.Errorf("unexpected payment status: %s", res.Status))
	}

	return &res, err
}
