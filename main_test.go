package main

import (
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"testing"
	"time"
)

type mockOpKiller struct {
	OpsKilled []Op
}

func (mok *mockOpKiller) Kill(op Op) error {
	mok.OpsKilled = append(mok.OpsKilled, op)
	return nil
}

type mockOpFinder struct {
	LastQuery   bson.M
	OpsReturned []Op
}

func (mof *mockOpFinder) Find(query bson.M) ([]Op, error) {
	mof.LastQuery = query
	return mof.OpsReturned, nil
}

func TestWhackAnOp(t *testing.T) {
	c := make(chan time.Time)
	ops := []Op{Op{
		ID:     1,
		Active: true,
		Op:     "query",
	}}
	query := bson.M{"some": "query"}
	finder := mockOpFinder{
		OpsReturned: ops,
	}
	killer := mockOpKiller{}
	whackanop := WhackAnOp{
		OpFinder: &finder,
		OpKiller: &killer,
		Query:    query,
		Tick:     c,
		Debug:    false,
	}
	done := make(chan bool)
	go func() {
		if err := whackanop.Run(); err != nil {
			t.Fatalf("unexpected error: %s", err)
		}
		done <- true
	}()
	c <- time.Now()
	close(c)
	<-done
	if !reflect.DeepEqual(finder.LastQuery, query) {
		t.Fatalf("expected query %#v, got %#v", query, finder.LastQuery)
	}
	if !reflect.DeepEqual(killer.OpsKilled, ops) {
		t.Fatalf("expected to kill %#v, got %#v", ops, killer.OpsKilled)
	}
}
