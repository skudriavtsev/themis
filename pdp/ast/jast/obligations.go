package jast

import (
	"encoding/json"

	"github.com/infobloxopen/themis/jparser"
	"github.com/infobloxopen/themis/pdp"
)

func (ctx context) unmarshalObligationItem(d *json.Decoder) (pdp.AttributeAssignmentExpression, error) {
	var (
		a pdp.Attribute
		e pdp.Expression
	)

	if err := jparser.UnmarshalObject(d, func(k string, d *json.Decoder) error {
		var (
			ok  bool
			err error
		)

		a, ok = ctx.attrs[k]
		if !ok {
			return newUnknownAttributeError(k)
		}

		if err = jparser.CheckObjectStart(d, "argument"); err == nil {
			e, err = ctx.unmarshalExpression(d)
			if err != nil {
				return bindError(err, k)
			}
		} else {
			return bindError(err, k)
		}
		return nil
	}, "obligation"); err != nil {
		return pdp.AttributeAssignmentExpression{}, err
	}

	return pdp.MakeAttributeAssignmentExpression(a, e), nil
}

func (ctx *context) unmarshalObligations(d *json.Decoder) ([]pdp.AttributeAssignmentExpression, error) {
	if err := jparser.CheckArrayStart(d, "obligations"); err != nil {
		return nil, err
	}

	var r []pdp.AttributeAssignmentExpression
	if err := jparser.UnmarshalObjectArray(d, func(idx int, d *json.Decoder) error {
		o, err := ctx.unmarshalObligationItem(d)
		if err != nil {
			return bindErrorf(err, "%d", idx)
		}

		r = append(r, o)

		return nil
	}, "obligations"); err != nil {
		return nil, err
	}

	return r, nil
}
