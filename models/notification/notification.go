package notification

// PendingNotification DB structure for pending notifications
type PendingNotification struct {
	ID        int
	Email     string
	Server    string
	Processor string
	Cores     int
	Threads   int
	Memory    int
	Storage   string
	Country   string
	Hardware  string
}
