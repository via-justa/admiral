// nolint:
package cmd

import (
	"reflect"
	"testing"
)

var inv = `{
    "_meta": {
        "hostvars": {
            "host1.domain.local": {
                "ansible_ssh_host": "1.1.1.1",
                "host_var1": {
                    "host_sub_var1": "host_sub_val1"
                }
            },
            "host2.domain.local": {
                "ansible_ssh_host": "2.2.2.2",
                "host_var2": "host_val2"
            },
            "host3.domain.local": {
                "ansible_ssh_host": "3.3.3.3",
                "host_var3": "host_val3"
            }
        }
    },
    "group1": {
        "hosts": [
            "host1.domain.local"
        ],
        "vars": {
            "group_var1": {
                "group_sub_var1": "group_sub_val1"
            }
        }
    },
    "group2": {
        "hosts": [
            "host2.domain.local"
        ],
        "vars": {
            "group_var2": "group_val2"
        }
    },
    "group3": {
        "hosts": [
            "host3.domain.local"
        ],
        "vars": {
            "group_var3": "group_val3"
        }
    },
    "group4": {
        "children": [
            "group3"
        ],
        "vars": {
            "group_var4": "group_val4"
        }
    },
    "group5": {
        "children": [
            "group4"
        ],
        "vars": {
            "group_var5": "group_val5"
        }
    }
}`

func Test_inventory(t *testing.T) {
	testDB := prepEnv()

	defer testDB.Close()

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
