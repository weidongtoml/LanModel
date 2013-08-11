// Copyright 2013 Weidong Liang. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package chinese_segmenter

import (
	"ngram_model"
	"strings"
	"testing"
)

const (
	unigram_path    = "../../data/model/unigram.dat"
	bigram_path     = "../../data/model/bigram.dat"
	cedict_path     = "../../data/dict/cedict_ts.u8.txt"
	cedict_key_type = TRADITION_CHINESE
)

func TestSegment(t *testing.T) {
	//cedict, err := LoadCEDict(cedict_path, cedict_key_type)
	//if err != nil {
	//	t.Fatalf("Failed to load CEDict[%s]: %s", cedict_path, err)
	//}
	model, err := ngram_model.LoadNGramModel(unigram_path, bigram_path)
	if err != nil {
		t.Fatalf("Failed to load model[%s,%s]: %s", unigram_path, bigram_path, err)
	}
	segmenter := NewSegmenter(nil, model)

	sentences := []string{
		"美 加緊 調查 襲擊 事件 宣佈 重新 開放 領空",
		"綜合 本 社 駐外 記者 報道 ， 在 美國 紐約 和 華盛頓 等 地 遭到 恐怖分子 襲擊 後 ， 美國 政府 正 加緊 進行 大 規模 的 調查 和 搜索 行動 ， 以 確定 製造 事件 的 恐怖 犯罪分子 。",
		"美國 總統 布什 12 日 與 世界 許多 國家 領導人 進行 廣泛 接觸 ， 努力 尋求 建立 一 個 反對 各 種 形式 的 恐怖主義 的 國際 聯盟 。",
		"布什 當天 分別 與 俄羅斯 、 法國 、 中國 、 英國 、 德國 和 加拿大 等 國 領導人 以及 聯合國 秘書長 安南 等 國際 組織 的 領導人 進行 了 電話 交談 。",
		"他 說 ， 建立 一 個 反對 各 種 形式 的 恐怖主義 的 國際 聯盟 是 美國 政府 目前 的 一 項 主要 任務 。",
		"當天 下午 ， 布什 視察 了 被 襲擊 的 五角大樓 。",
		"此前 ， 他 在 白宮 舉行 安全 會議 並 聽取 了 情報 部門 的 最 新 匯報 。",
		"美國 政府 在 外國 安全 部門 的 合作 下 ， 正 加緊 進行 搜索 恐怖分子 的 行動 。",
		"美國 司法部長 阿什克羅夫特 12 日 在 聯邦 調查局 舉行 的 記者 招待會 上 說 ， 一些 劫機 嫌疑犯 已 被 確認 ， 他們 曾 在 美國 接受 過 駕駛 飛機 的 訓練 。",
		"美國 聯邦 調查局 局長 米勒 也 在 記者 招待會 上 說 ， 美國 政府 正在 對 一些 劫機 嫌疑犯 展開 取證 工作 。",
		"他 說 ， 4000 名 聯邦 調查局 專家 以及 3000 名 支援 人員 正在 展開 這 場 美國 歷史 上 最 大 的 調查 和 搜索 行動 ， 另有 400 多 名 在 聯邦 調查局 刑事 實驗室 工作 的 人員 被 派 往 出事 地點 。",
		"據 報道 ， 聯邦 特工 在 加利福尼亞州 和 波士頓 等 地 展開 了 搜查 行動 。",
		"此外 ， 德國 警方 也 在 漢堡 搜查 了 一 處 疑犯 曾經 居住 過 的 住所 。",
		"美國 運輸部長 諾曼．峰田 13 日 宣佈 ， 從 美國 東部 時間 上午 11 時 起 ， 美國 重新 開放 領空 ， 允許 商業 和 私人 飛機 恢復 飛行 。",
		"峰田 說 ， 鑑於 航班 恢復 較 慢 和 機場 安全 措施 加強 ， 乘客 應 事先 向 航空 公司 諮詢 有關 航班 時間 和 服務 事宜 ， 並 留 出 足夠 的 時間 以備 安全 檢查 。",
	}

	for _, s := range sentences {
		sent := strings.Replace(s, " ", "", -1)
		terms := strings.Split(s, " ")
		result := segmenter.Segment(sent)
		is_eqv := len(result) == len(terms)
		for i, r := range result {
			if r != terms[i] {
				is_eqv = false
				break
			}
		}
		if !is_eqv {
			t.Errorf("Expected result to be %v but got %v", terms, result)
		}
	}
}