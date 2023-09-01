#! /bin/bash
# Author:ZhuYunFei
# 监听主程的 SIGTERM 信号并透传给当前服务
# 主要用于触发服务的优雅退出流程
# 愿世界和平

./yourshinesapi & pid="$!" # 启动服务并记录 pid

sigterm_handle() {
  echo "[INFO] 通知服务开始退出"
  kill -TERM $pid  # 传递 SIGTERM 给当前服务
  wait $pid
}

trap sigterm_handle TERM # 捕获 SIGTERM 信号并回调 handle_sigterm 函数
echo "[INFO] 服务启动成功，PID为：$pid"

wait # 等待回调执行完，主进程再退出
