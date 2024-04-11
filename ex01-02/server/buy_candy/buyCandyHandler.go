package buy_candy
/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

unsigned int i;
unsigned int argscharcount = 0;

char *ask_cow(char phrase[]) {
  int phrase_len = strlen(phrase);
  char *buf = (char *)malloc(sizeof(char) * (160 + (phrase_len + 2) * 3));
  strcpy(buf, " ");

  for (i = 0; i < phrase_len + 2; ++i) {
    strcat(buf, "_");
  }

  strcat(buf, "\n< ");
  strcat(buf, phrase);
  strcat(buf, " ");
  strcat(buf, ">\n ");

  for (i = 0; i < phrase_len + 2; ++i) {
    strcat(buf, "-");
  }
  strcat(buf, "\n");
  strcat(buf, "        \\   ^__^\n");
  strcat(buf, "         \\  (oo)\\_______\n");
  strcat(buf, "            (__)\\       )\\/\\\n");
  strcat(buf, "                ||----w |\n");
  strcat(buf, "                ||     ||\n");
  return buf;
}
*/
import "C"

import (
	"fmt"
	"operations"
	"github.com/go-openapi/runtime/middleware"
)

var prices = map[string]int64 {
	"CE": 10,
	"AA": 15,
	"NT": 17,
	"DE": 21,
	"YR": 23,
}


func BuyCandyHandler(b operations.BuyCandyParams) middleware.Responder {
	candy_price, ok := prices[*b.Order.CandyType]
	if !ok || *b.Order.Money < 0 || *b.Order.CandyCount < 0 {
		return operations.NewBuyCandyBadRequest().WithPayload(
			&operations.BuyCandyBadRequestBody{
				Error: AskCow("some error in input data"),
			},
		)
	}

	var response middleware.Responder
	price := *b.Order.CandyCount * candy_price
	if price > 0 && price <= *b.Order.Money {
		response = operations.NewBuyCandyCreated().WithPayload(
			&operations.BuyCandyCreatedBody{
				Thanks: AskCow("Thank you! Your change is " +
					fmt.Sprintln(*b.Order.Money - price)),
				Change: (*b.Order.Money - price),
			},
		)
	} else if price > 0 && price > *b.Order.Money {
		response = operations.NewBuyCandyPaymentRequired().WithPayload(
			&operations.BuyCandyPaymentRequiredBody{
				Error: AskCow(fmt.Sprintf("You need %d more money", price - *b.Order.Money, )),
			},
		)
	}
	return response
}

func AskCow(str string) string {
	s := C.CString(str)
	out := C.ask_cow(s)
	res := C.GoString(out)
	return res
}