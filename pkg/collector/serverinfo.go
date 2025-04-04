package collector

import (
	"log"

	"github.com/hikhvar/ts3exporter/pkg/serverquery"
	"github.com/prometheus/client_golang/prometheus"
)

const serverInfoSubsystem = "serverinfo"

var serverInfoLabels = []string{virtualServerLabel}

type ServerInfo struct {
	executor        serverquery.Executor
	internalMetrics *ExporterMetrics

	ClientsOnline             *prometheus.Desc
	QueryClientsOnline        *prometheus.Desc
	Online                    *prometheus.Desc
	MaxClients                *prometheus.Desc
	Uptime                    *prometheus.Desc
	ChannelsOnline            *prometheus.Desc
	MaxDownloadTotalBandwidth *prometheus.Desc
	MaxUploadTotalBandwidth   *prometheus.Desc
	ClientsConnections        *prometheus.Desc
	QueryClientsConnections   *prometheus.Desc

	FileTransferBytesSentTotal     *prometheus.Desc
	FileTransferBytesReceivedTotal *prometheus.Desc

	ControlBytesSentTotal     *prometheus.Desc
	ControlBytesReceivedTotal *prometheus.Desc

	SpeechBytesSentTotal     *prometheus.Desc
	SpeechBytesReceivedTotal *prometheus.Desc

	KeepAliveBytesSentTotal     *prometheus.Desc
	KeepAliveBytesReceivedTotal *prometheus.Desc

	BytesSentTotal     *prometheus.Desc
	BytesReceivedTotal *prometheus.Desc

	ControlPacketLoss   *prometheus.Desc
	SpeechPacketLoss    *prometheus.Desc
	KeepAlivePacketLoss *prometheus.Desc
	TotalPacketLoss     *prometheus.Desc

	Ping *prometheus.Desc
}

func NewServerInfo(executor serverquery.Executor, internalMetrics *ExporterMetrics) *ServerInfo {
	return &ServerInfo{
		executor:                       executor,
		internalMetrics:                internalMetrics,
		ClientsOnline:                  prometheus.NewDesc(fqdn(serverInfoSubsystem, "clients_online"), "number of currently online clients", serverInfoLabels, nil),
		QueryClientsOnline:             prometheus.NewDesc(fqdn(serverInfoSubsystem, "query_clients_online"), "number of currently online query clients", serverInfoLabels, nil),
		Online:                         prometheus.NewDesc(fqdn(serverInfoSubsystem, "online"), "is the virtual server online", serverInfoLabels, nil),
		MaxClients:                     prometheus.NewDesc(fqdn(serverInfoSubsystem, "max_clients"), "maximal number of allowed clients", serverInfoLabels, nil),
		Uptime:                         prometheus.NewDesc(fqdn(serverInfoSubsystem, "uptime"), "uptime of the virtual server", serverInfoLabels, nil),
		ChannelsOnline:                 prometheus.NewDesc(fqdn(serverInfoSubsystem, "channels_online"), "number of online channels", serverInfoLabels, nil),
		MaxDownloadTotalBandwidth:      prometheus.NewDesc(fqdn(serverInfoSubsystem, "download_bandwidth_bytes_max"), "maximal bandwidth available for downloads", serverInfoLabels, nil),
		MaxUploadTotalBandwidth:        prometheus.NewDesc(fqdn(serverInfoSubsystem, "upload_bandwidth_bytes_max"), "maximal bandwidth available for uploads", serverInfoLabels, nil),
		ClientsConnections:             prometheus.NewDesc(fqdn(serverInfoSubsystem, "client_connections"), "currently established client connections", serverInfoLabels, nil),
		QueryClientsConnections:        prometheus.NewDesc(fqdn(serverInfoSubsystem, "query_client_connections"), "currently established query client connections", serverInfoLabels, nil),
		FileTransferBytesSentTotal:     prometheus.NewDesc(fqdn(serverInfoSubsystem, "file_transfer_bytes_sent_total"), "total sent bytes for file transfers", serverInfoLabels, nil),
		FileTransferBytesReceivedTotal: prometheus.NewDesc(fqdn(serverInfoSubsystem, "file_transfer_bytes_received_total"), "total received bytes for file transfers", serverInfoLabels, nil),
		ControlBytesSentTotal:          prometheus.NewDesc(fqdn(serverInfoSubsystem, "control_bytes_sent_total"), "total sent bytes for control traffic", serverInfoLabels, nil),
		ControlBytesReceivedTotal:      prometheus.NewDesc(fqdn(serverInfoSubsystem, "control_bytes_received_total"), "total received bytes for control traffic", serverInfoLabels, nil),
		SpeechBytesSentTotal:           prometheus.NewDesc(fqdn(serverInfoSubsystem, "speech_bytes_sent_total"), "total sent bytes for speech traffic", serverInfoLabels, nil),
		SpeechBytesReceivedTotal:       prometheus.NewDesc(fqdn(serverInfoSubsystem, "speech_bytes_received_total"), "total received bytes for speech traffic", serverInfoLabels, nil),
		KeepAliveBytesSentTotal:        prometheus.NewDesc(fqdn(serverInfoSubsystem, "keepalive_bytes_sent_total"), "total send bytes for keepalive traffic", serverInfoLabels, nil),
		KeepAliveBytesReceivedTotal:    prometheus.NewDesc(fqdn(serverInfoSubsystem, "keepalive_bytes_received_total"), "total received bytes for keepalive traffic", serverInfoLabels, nil),
		BytesSentTotal:                 prometheus.NewDesc(fqdn(serverInfoSubsystem, "bytes_send_total"), "total send bytes", serverInfoLabels, nil),
		BytesReceivedTotal:             prometheus.NewDesc(fqdn(serverInfoSubsystem, "bytes_received_total"), "total received bytes", serverInfoLabels, nil),
		ControlPacketLoss:              prometheus.NewDesc(fqdn(serverInfoSubsystem, "control_packet_loss"), "packet loss in control traffic", serverInfoLabels, nil),
		SpeechPacketLoss:               prometheus.NewDesc(fqdn(serverInfoSubsystem, "speech_packet_loss"), "packet loss in speech traffic", serverInfoLabels, nil),
		KeepAlivePacketLoss:            prometheus.NewDesc(fqdn(serverInfoSubsystem, "keepalive_packet_loss"), "packet loss in keepalive traffic", serverInfoLabels, nil),
		TotalPacketLoss:                prometheus.NewDesc(fqdn(serverInfoSubsystem, "total_packet_loss"), "packet loss in total traffic", serverInfoLabels, nil),
		Ping:                           prometheus.NewDesc(fqdn(serverInfoSubsystem, "ping"), "average ping of all connected clients", serverInfoLabels, nil),
	}

}

func (s *ServerInfo) Describe(c chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(s, c)
}

func (s *ServerInfo) Collect(c chan<- prometheus.Metric) {
	vServerView := serverquery.NewVirtualServer(s.executor)
	if err := vServerView.Refresh(); err != nil {
		s.internalMetrics.RefreshError(serverInfoSubsystem)
		log.Printf("failed to refresh server info view: %v", err)
	}
	for _, vs := range vServerView.All() {
		c <- prometheus.MustNewConstMetric(s.ClientsOnline, prometheus.GaugeValue, float64(vs.ClientsOnline), vs.Name)
		c <- prometheus.MustNewConstMetric(s.QueryClientsOnline, prometheus.GaugeValue, float64(vs.QueryClientsOnline), vs.Name)
		c <- prometheus.MustNewConstMetric(s.Online, prometheus.GaugeValue, online(vs.Status), vs.Name)
		c <- prometheus.MustNewConstMetric(s.MaxClients, prometheus.GaugeValue, float64(vs.MaxClients), vs.Name)
		c <- prometheus.MustNewConstMetric(s.Uptime, prometheus.CounterValue, float64(vs.Uptime), vs.Name)
		c <- prometheus.MustNewConstMetric(s.ChannelsOnline, prometheus.GaugeValue, float64(vs.ChannelsOnline), vs.Name)
		c <- prometheus.MustNewConstMetric(s.MaxDownloadTotalBandwidth, prometheus.GaugeValue, float64(vs.MaxDownloadTotalBandwidth), vs.Name)
		c <- prometheus.MustNewConstMetric(s.MaxUploadTotalBandwidth, prometheus.GaugeValue, float64(vs.MaxUploadTotalBandwidth), vs.Name)
		c <- prometheus.MustNewConstMetric(s.ClientsConnections, prometheus.GaugeValue, float64(vs.ClientsConnections), vs.Name)
		c <- prometheus.MustNewConstMetric(s.QueryClientsConnections, prometheus.GaugeValue, float64(vs.QueryClientsConnections), vs.Name)
		c <- prometheus.MustNewConstMetric(s.FileTransferBytesSentTotal, prometheus.CounterValue, float64(vs.FileTransferBytesSentTotal), vs.Name)
		c <- prometheus.MustNewConstMetric(s.FileTransferBytesReceivedTotal, prometheus.CounterValue, float64(vs.FileTransferBytesReceivedTotal), vs.Name)
		c <- prometheus.MustNewConstMetric(s.ControlBytesSentTotal, prometheus.CounterValue, float64(vs.ControlBytesSentTotal), vs.Name)
		c <- prometheus.MustNewConstMetric(s.ControlBytesReceivedTotal, prometheus.CounterValue, float64(vs.ControlBytesReceivedTotal), vs.Name)
		c <- prometheus.MustNewConstMetric(s.SpeechBytesSentTotal, prometheus.CounterValue, float64(vs.SpeechBytesSentTotal), vs.Name)
		c <- prometheus.MustNewConstMetric(s.SpeechBytesReceivedTotal, prometheus.CounterValue, float64(vs.SpeechBytesReceivedTotal), vs.Name)
		c <- prometheus.MustNewConstMetric(s.KeepAliveBytesSentTotal, prometheus.CounterValue, float64(vs.KeepAliveBytesSentTotal), vs.Name)
		c <- prometheus.MustNewConstMetric(s.KeepAliveBytesReceivedTotal, prometheus.CounterValue, float64(vs.KeepAliveBytesReceivedTotal), vs.Name)
		c <- prometheus.MustNewConstMetric(s.BytesSentTotal, prometheus.CounterValue, float64(vs.BytesSentTotal), vs.Name)
		c <- prometheus.MustNewConstMetric(s.BytesReceivedTotal, prometheus.CounterValue, float64(vs.BytesReceivedTotal), vs.Name)
		c <- prometheus.MustNewConstMetric(s.ControlPacketLoss, prometheus.CounterValue, vs.ControlPacketLoss, vs.Name)
		c <- prometheus.MustNewConstMetric(s.SpeechPacketLoss, prometheus.CounterValue, vs.SpeechPacketLoss, vs.Name)
		c <- prometheus.MustNewConstMetric(s.KeepAlivePacketLoss, prometheus.CounterValue, vs.KeepAlivePacketLoss, vs.Name)
		c <- prometheus.MustNewConstMetric(s.TotalPacketLoss, prometheus.CounterValue, vs.ControlPacketLoss, vs.Name)
		c <- prometheus.MustNewConstMetric(s.Ping, prometheus.CounterValue, vs.Ping, vs.Name)
	}
}

func online(status string) float64 {
	if status == "online" {
		return 1.0
	}
	return 0.0
}
