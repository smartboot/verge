#!/bin/sh
dir=$(dirname "$0")

LOG_PATH=./logs/ibms.log


cat > /lib/systemd/system/verge.service <<EOF
[Unit]
Description=verge
After=network.target

[Service]
Type=simple
WorkingDirectory=${dir}
ExecStart=${dir}/verge
Restart=always
RestartSec=5s
TimeoutSec=5s
Environment=DRIVERBOX_LOG_PATH=${LOG_PATH}
Environment=EXPORT_IBUILDING_LICENSE=${LICENSE}
Environment=EXPORT_IBUILDING_AGENT_BROKER=${AGENT_BROKER}
Environment=EXPORT_IBUILDING_ENABLED=${EXPORT_IBUILDING_ENABLED}
Environment=EXPORT_IBUILDING_LOW_FLOW=${EXPORT_IBUILDING_LOW_FLOW}
Environment=EXPORT_IBUILDING_TIMER_REPORT_PERIOD=${EXPORT_IBUILDING_TIMER_REPORT_PERIOD}
Environment=EXPORT_COMPUTING_ENABLED=${EXPORT_COMPUTING_ENABLED}
Environment=EXPORT_IBMS_ENABLED=${EXPORT_IBMS_ENABLED}
Environment=DRIVERBOX_VIRTUAL=${DRIVERBOX_VIRTUAL}
Environment=DRIVERBOX_COMPUTING_VIRTUAL=${DRIVERBOX_COMPUTING_VIRTUAL}
Environment=PPROF_ENABLED=${PPROF_ENABLED}
Environment=ENV_AUTO_DISCOVERY=${ENV_AUTO_DISCOVERY}
Environment=UPGRADE_SAVE_PATH=${UPGRADE_SAVE_PATH}
Environment=HARDWARE_RUNTIME=M0

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload
systemctl start verge
systemctl enable verge
echo 'systemctl status verge'