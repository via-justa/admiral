// nolint:
package cmd

import (
	"reflect"
	"testing"
)

var inv = `{
    "_meta": {
        "hostvars": {
            "host1.local": {
                "ansible_ssh_host": "1.1.1.1",
                "host_var1": {
                    "host_sub_var1": "host_sub_val1"
                }
            },
            "host2.local": {
                "ansible_ssh_host": "2.2.2.2",
                "var2": "val2"
            },
            "host3.local": {
                "ansible_ssh_host": "3.3.3.3",
                "var3": "val3"
            }
        }
    },
    "group1": {
        "hosts": [
            "host1.local"
        ],
        "vars": {
            "group_var1": {
                "group_sub_var1": "group_sub_val1"
            }
        }
    },
    "group2": {
        "hosts": [
            "host2.local"
        ],
        "vars": {
            "var2": "val2"
        }
    },
    "group3": {
        "hosts": [
            "host3.local"
        ],
        "vars": {
            "var3": "val3"
        }
    },
    "group4": {
        "children": [
            "group3"
        ],
        "vars": {
            "var4": "val4"
        }
    },
    "group5": {
        "children": [
            "group4"
        ],
        "vars": {
            "var5": "val5"
        }
    }
}`

func Test_inventory(t *testing.T) {
	db = dbMock{}

	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		{
			name:    "Generate inventory",
			want:    []byte(inv),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inventory()
			if (err != nil) != tt.wantErr {
				t.Errorf("inventory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inventory() = %s, want %s", got, tt.want)
			}
		})
	}
}
