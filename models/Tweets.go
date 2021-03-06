package models

import (
	"fmt"
	"strings"
	"time"
)

type Tweet struct {
	CreatedAt string `json:"created_at"`
	Id int `json:"id"`
	Text string `json:"text"`
	Origin string `json:"origin"`
	Tag string `json:"tag"`
}

// ChangeDateFormat 날짜 포맷을 변경한다
// Before "Thu Apr 01 01:54:29 +0000 2021" : time.RubyDate
// After "2021-04-01T10:54:29+09:00" : time.RFC3339
func (w *Tweet) ChangeDateFormat() error {
	if w == nil {
		return fmt.Errorf("receiver is nil")
	}

	// created_at 필드를 파싱하여 Time 객체로 생성
	t, err := time.Parse(time.RubyDate, w.CreatedAt)
	if err != nil {
		return err
	}

	// 타임존 로케이션 설정
	loc, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return err
	}

	// 타임존 변경 및 RFC3339 형식으로 날짜 형식 변경
	tToKst := t.In(loc).Format(time.RFC3339)

	// 변경한 created_at으로 변경
	w.CreatedAt = tToKst

	return nil
}

// SetOrigin 데이터의 출처를 설정한다
func (w *Tweet) SetOrigin() error {
	if w == nil {
		return fmt.Errorf("receiver is nil")
	}

	w.Origin = "twitter"
	return nil
}

// SetTag 태그를 설정한다
// 태그는 프로젝트에서의 메시지 카테고리를 의미한다
func (w *Tweet) SetTag() error{
	if w== nil {
		return fmt.Errorf("receiver is nil")
	}
	w.Tag = "sns"
	return nil
}

// TrimText Text 데이터의 공백, 개행을 제거한다
func (w *Tweet) TrimText() {
	if w == nil {
		return
	}
	txt := w.Text
	txt = strings.Replace(txt, "\n", "", -1)
	w.Text = strings.Trim(txt, " ")
}

// IsRetweet 리트윗이라면 무시한다
// 트윗 내용이 'RT @' 으로 시작하면 true를 반환한다
func (w *Tweet) IsRetweet() bool {
	return strings.HasPrefix(w.Text, "RT @")
}