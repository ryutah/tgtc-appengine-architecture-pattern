# TGTC AppEngineアーキテクチャパターン @ryutah
## サンプルで利用するGCPリソースの準備
```console
$ ./resource/provisioning.sh [DEPLOY_NAME]
```

## 各サンプルアプリケーション実行前の設定
```console
$ export GOOGLE_CLOUD_PROJECT=[YOUR_PROJECT_ID]
$ export GCS_BUCKET=[YOUR_PROJECT_ID]-[DEPLOY_NAME]
```
