# VRC Lottery System

[![License](https://img.shields.io/badge/License-AGPL-blue.svg)](LICENSE)

VRChat内の小規模(700ユーザー以下)イベントでの抽選を想定した抽選システム  
> [!WARNING]
> 開発中なので現時点で使用できません

## Build
```bash
CGO_ENABLED=0 GOOS=linux go build -a -trimpath -ldflags="-s -w -X main.Version={{Version}}"  -o VRCLotterySystem
```

## TODO
- [ ] Unit Test
- [ ] 管理用Discord botの実装
- [ ] フロントエンドの実装
- [ ] 基本ロジックの見直し

## Project Structure
```
├── frontend                      # フロントエンドのソースコード（未実装）
├── main.go                       
├── config/
│   └── config.go                 # 設定ファイル(JSON)の定義・読み込み・バリデーション
├── log/
│   └── log.go                    # グローバルLogger
├── internal/                     
│   ├── global/
│   │   └── validator/
│   │       └── validator.go      # グローバルvalidator
│   ├── http/                     # HTTP サーバ (Gin)
│   │   ├── http.go               # サーバ初期化・ミドルウェア登録
│   │   ├── route.go              # ルーティング定義
│   │   ├── handle.go             # ハンドラ
│   │   ├── middlware.go          # リクエストログ等のミドルウェア
│   │   └── resp.go               # 共通レスポンス構造体
│   ├── provider/                 # ロジックの実装及び上層への提供
│   │   ├── provider.go
│   │   ├── drawing.go            # 抽選処理
│   │   ├── inputs.go             # 応募フォームデータの処理
│   │   ├── results.go            # 抽選結果
│   │   ├── blacklist.go          # ブラックリスト
│   │   └── stafflist.go          # スタッフリスト
│   ├── data/
│   │   └── data.go               # 永続データの read/write（JSON, 排他制御付き）
│   ├── discord/
│   │   └── discord.go            # Discord bot連携（未実装）
│   └── task/
│       └── task.go               # 定期実行タスク（抽選バッチ）
└── pkg/                          
    └── eth/
        ├── eth.go                # Ethereumのブロック高度による乱数生成ロジック
        └── eth_test.go
```


## ライセンス
本プロジェクトは [GNUアフェロ一般公衆ライセンス](https://gpl.mhatta.org/agpl.ja.html) を基づき発行しております、
被配布者(ユーザー)には使用の自由・二次開発の自由・二次配布（販売含む）の自由などの自由が保証されます。
ただし二次配布の場合必ず同じ[GNUアフェロ一般公衆ライセンス](https://gpl.mhatta.org/agpl.ja.html) で発行するよう義務付けられます、
そして二次被配布者には本プロジェクトが被配布者に与えた自由と同じ自由が与えられます、
それらの自由を制限・干渉するあらゆる行為は一切認められません。