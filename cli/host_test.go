package cli

import (
	"reflect"
	"testing"

	"github.com/via-justa/admiral/datastructs"
)

func TestCreateHost(t *testing.T) {
	db = dbMock{}

	type args struct {
		host datastructs.Host
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Insert host",
			args: args{
				host: datastructs.Host{
					ID:        2,
					Hostname:  "host3",
					Host:      "3.3.3.3",
					Variables: "{\"var3\": \"val3\"}",
					Enabled:   true,
					Monitored: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Insert Existing host without change",
			args: args{
				host: datastructs.Host{
					ID:        1,
					Hostname:  "host1",
					Host:      "1.1.1.1",
					Variables: "{\"var1\": \"val1\"}",
					Enabled:   true,
					Monitored: true,
				},
			},
			wantErr: true,
		},
		{
			name: "Change Existing host",
			args: args{
				host: datastructs.Host{
					ID:        1,
					Hostname:  "host1",
					Host:      "1.1.1.1",
					Variables: "{\"var1\": \"val1\"}",
					Enabled:   true,
					Monitored: false,
				},
			},
			wantErr: false,
		},
		{
			name: "Missing host HostName",
			args: args{
				host: datastructs.Host{Hostname: ""},
			},
			wantErr: true,
		},
		{
			name: "Missing host Host",
			args: args{
				host: datastructs.Host{Host: ""},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateHost(tt.args.host); (err != nil) != tt.wantErr {
				t.Errorf("CreateHost() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestViewHostByHostname(t *testing.T) {
	db = dbMock{}

	type args struct {
		hostname string
	}
	tests := []struct {
		name     string
		args     args
		wantHost datastructs.Host
		wantErr  bool
	}{
		{
			name: "View host by HostName",
			args: args{
				hostname: "host1",
			},
			wantHost: datastructs.Host{
				ID:        1,
				Hostname:  "host1",
				Host:      "1.1.1.1",
				Domain:    "local",
				Variables: "{\"var1\": \"val1\"}",
				Enabled:   true,
				Monitored: true,
				Groups:    []string{"group1", "group2", "group3"},
			},
			wantErr: false,
		},
		{
			name: "None existing host",
			args: args{
				hostname: "host2",
			},
			wantHost: datastructs.Host{},
			wantErr:  true,
		},
		{
			name: "Empty host HostName",
			args: args{
				hostname: "",
			},
			wantHost: datastructs.Host{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, err := ViewHostByHostname(tt.args.hostname)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewHostByHostname() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHost, tt.wantHost) {
				t.Errorf("ViewHostByHostname() = %v, want %v", gotHost, tt.wantHost)
			}
		})
	}
}

func TestViewHostByIP(t *testing.T) {
	db = dbMock{}

	type args struct {
		ip string
	}
	tests := []struct {
		name     string
		args     args
		wantHost datastructs.Host
		wantErr  bool
	}{
		{
			name: "View host by IP",
			args: args{
				ip: "1.1.1.1",
			},
			wantHost: datastructs.Host{
				ID:        1,
				Hostname:  "host1",
				Host:      "1.1.1.1",
				Domain:    "local",
				Variables: "{\"var1\": \"val1\"}",
				Enabled:   true,
				Monitored: true,
				Groups:    []string{"group1", "group2", "group3"},
			},
			wantErr: false,
		},
		{
			name: "None existing host",
			args: args{
				ip: "2.2.2.2",
			},
			wantHost: datastructs.Host{},
			wantErr:  true,
		},
		{
			name: "Empty host Host",
			args: args{
				ip: "",
			},
			wantHost: datastructs.Host{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, err := ViewHostByIP(tt.args.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewHostByIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHost, tt.wantHost) {
				t.Errorf("ViewHostByIP() = %v, want %v", gotHost, tt.wantHost)
			}
		})
	}
}

func TestViewHostByID(t *testing.T) {
	db = dbMock{}

	type args struct {
		id int
	}
	tests := []struct {
		name     string
		args     args
		wantHost datastructs.Host
		wantErr  bool
	}{
		{
			name: "View host by IP",
			args: args{
				id: 1,
			},
			wantHost: datastructs.Host{
				ID:        1,
				Hostname:  "host1",
				Host:      "1.1.1.1",
				Domain:    "local",
				Variables: "{\"var1\": \"val1\"}",
				Enabled:   true,
				Monitored: true,
				Groups:    []string{"group1", "group2", "group3"},
			},
			wantErr: false,
		},
		{
			name: "None existing host",
			args: args{
				id: 2,
			},
			wantHost: datastructs.Host{},
			wantErr:  true,
		},
		{
			name: "Empty host id",
			args: args{
				id: 0,
			},
			wantHost: datastructs.Host{},
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHost, err := ViewHostByID(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("ViewHostByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHost, tt.wantHost) {
				t.Errorf("ViewHostByID() = %v, want %v", gotHost, tt.wantHost)
			}
		})
	}
}

func TestListHosts(t *testing.T) {
	db = dbMock{}

	tests := []struct {
		name      string
		wantHosts []datastructs.Host
		wantErr   bool
	}{
		{
			name: "List hosts",
			wantHosts: []datastructs.Host{
				datastructs.Host{
					ID:        1,
					Hostname:  "host1",
					Host:      "1.1.1.1",
					Domain:    "local",
					Variables: "{\"var1\": \"val1\"}",
					Enabled:   true,
					Monitored: true,
					Groups:    []string{"group1", "group2", "group3"},
				},
				datastructs.Host{
					ID:        2,
					Hostname:  "host2",
					Host:      "2.2.2.2",
					Domain:    "local",
					Variables: "{\"var2\": \"val2\"}",
					Enabled:   true,
					Monitored: false,
					Groups:    []string{},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHosts, err := ListHosts()
			if (err != nil) != tt.wantErr {
				t.Errorf("ListHosts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotHosts, tt.wantHosts) {
				t.Errorf("ListHosts() = %v, want %v", gotHosts, tt.wantHosts)
			}
		})
	}
}

func TestDeleteHost(t *testing.T) {
	db = dbMock{}

	type args struct {
		host datastructs.Host
	}
	tests := []struct {
		name         string
		args         args
		wantAffected int64
		wantErr      bool
	}{
		{
			name: "Delete host",
			args: args{
				host: datastructs.Host{ID: 1},
			},
			wantAffected: 1,
			wantErr:      false,
		},
		{
			name: "Delete non-existing host",
			args: args{
				host: datastructs.Host{ID: 3},
			},
			wantAffected: 0,
			wantErr:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotAffected, err := DeleteHost(tt.args.host)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteHost() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotAffected != tt.wantAffected {
				t.Errorf("DeleteHost() = %v, want %v", gotAffected, tt.wantAffected)
			}
		})
	}
}
