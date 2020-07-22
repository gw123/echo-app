package components

import (
	"crypto/rsa"
	"testing"
	"time"
)

func TestJwsHelper_CreateTokenAndParseToken(t *testing.T) {
	type fields struct {
		opt        JwsHelperOpt
		publicKey  *rsa.PublicKey
		privateKey *rsa.PrivateKey
	}
	type args struct {
		userId  int64
		payload string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		want        string
		wantErr     bool
		wantTimeout bool
		timeout     int
	}{
		{
			name: "001",
			fields: fields{
				opt: JwsHelperOpt{
					Audience:       "xyt",
					Issuer:         "gw123",
					Timeout:        3600 * 24,
					PublicKeyPath:  "./test/sample_key.pub",
					PrivateKeyPath: "./test/sample_key",
					HashIdsSalt:    "123456",
				},
			},
			args: args{
				userId:  123,
				payload: "hello world",
			},
			//want:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ4eXQiLCJleHAiOjE1ODk0NjEwMjEsImp0aSI6IjEyMyIsImlhdCI6MTU4OTM3NDYyMSwiaXNzIjoiZ3cxMjMiLCJzdWIiOiJoZWxsbyB3b3JsZCJ9.M75zks254kTAH_h1CMEPmEZRgpl8OS_NM0SB6Xh1iNzKjORjgOZ6EvrdeROYi1HHhQBlePZFyqaOtJL2WoE7h4KuUVAOWA5VuR5xQQyAIDzxIr0qK_QFy0vS851vPBnqrrMN2rd04ATAN0_iFzgMXs5vqLjtLLmLBzYBMkuvxFsRVHXWgmBXqOxHPAj3wRWM73unGdaEPo9v_-dDkXj42XBUDOa-9-YVrfz5nNDObPDoaPcyYL032NVmA_IUC9xw2FUV63KB_rk-SupPbjK2jiQ6viK51NIBf4l_d1hxt8fIEsBdGQUNiyCsJhhevt_flIccIJ5Y1mTaF07FfgbTvA",
			wantErr: false,
		},
		{
			name: "002",
			fields: fields{
				opt: JwsHelperOpt{
					Audience:       "xyt",
					Issuer:         "gw123",
					Timeout:        3600 * 24,
					PublicKeyPath:  "./test/sample_key.pub",
					PrivateKeyPath: "./test/sample_key",
					HashIdsSalt:    "",
				},
			},
			args: args{
				userId:  123,
				payload: "hello world",
			},
			//want:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ4eXQiLCJleHAiOjE1ODk0NjEwMjEsImp0aSI6IjEyMyIsImlhdCI6MTU4OTM3NDYyMSwiaXNzIjoiZ3cxMjMiLCJzdWIiOiJoZWxsbyB3b3JsZCJ9.M75zks254kTAH_h1CMEPmEZRgpl8OS_NM0SB6Xh1iNzKjORjgOZ6EvrdeROYi1HHhQBlePZFyqaOtJL2WoE7h4KuUVAOWA5VuR5xQQyAIDzxIr0qK_QFy0vS851vPBnqrrMN2rd04ATAN0_iFzgMXs5vqLjtLLmLBzYBMkuvxFsRVHXWgmBXqOxHPAj3wRWM73unGdaEPo9v_-dDkXj42XBUDOa-9-YVrfz5nNDObPDoaPcyYL032NVmA_IUC9xw2FUV63KB_rk-SupPbjK2jiQ6viK51NIBf4l_d1hxt8fIEsBdGQUNiyCsJhhevt_flIccIJ5Y1mTaF07FfgbTvA",
			wantErr: false,
		},
		{
			name: "003",
			fields: fields{
				opt: JwsHelperOpt{
					Audience:       "xyt",
					Issuer:         "gw123",
					Timeout:        60,
					PublicKeyPath:  "./test/sample_key.pub",
					PrivateKeyPath: "./test/sample_key",
					HashIdsSalt:    "",
				},
			},
			args: args{
				userId:  123,
				payload: "hello world",
			},
			//want:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ4eXQiLCJleHAiOjE1ODk0NjEwMjEsImp0aSI6IjEyMyIsImlhdCI6MTU4OTM3NDYyMSwiaXNzIjoiZ3cxMjMiLCJzdWIiOiJoZWxsbyB3b3JsZCJ9.M75zks254kTAH_h1CMEPmEZRgpl8OS_NM0SB6Xh1iNzKjORjgOZ6EvrdeROYi1HHhQBlePZFyqaOtJL2WoE7h4KuUVAOWA5VuR5xQQyAIDzxIr0qK_QFy0vS851vPBnqrrMN2rd04ATAN0_iFzgMXs5vqLjtLLmLBzYBMkuvxFsRVHXWgmBXqOxHPAj3wRWM73unGdaEPo9v_-dDkXj42XBUDOa-9-YVrfz5nNDObPDoaPcyYL032NVmA_IUC9xw2FUV63KB_rk-SupPbjK2jiQ6viK51NIBf4l_d1hxt8fIEsBdGQUNiyCsJhhevt_flIccIJ5Y1mTaF07FfgbTvA",
			wantErr: false,
		},
		{
			name: "004",
			fields: fields{
				opt: JwsHelperOpt{
					Audience:       "xyt",
					Issuer:         "gw123",
					Timeout:        1,
					PublicKeyPath:  "./test/sample_key.pub",
					PrivateKeyPath: "./test/sample_key",
					HashIdsSalt:    "",
				},
			},
			args: args{
				userId:  123,
				payload: "hello world",
			},
			//want:    "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhdWQiOiJ4eXQiLCJleHAiOjE1ODk0NjEwMjEsImp0aSI6IjEyMyIsImlhdCI6MTU4OTM3NDYyMSwiaXNzIjoiZ3cxMjMiLCJzdWIiOiJoZWxsbyB3b3JsZCJ9.M75zks254kTAH_h1CMEPmEZRgpl8OS_NM0SB6Xh1iNzKjORjgOZ6EvrdeROYi1HHhQBlePZFyqaOtJL2WoE7h4KuUVAOWA5VuR5xQQyAIDzxIr0qK_QFy0vS851vPBnqrrMN2rd04ATAN0_iFzgMXs5vqLjtLLmLBzYBMkuvxFsRVHXWgmBXqOxHPAj3wRWM73unGdaEPo9v_-dDkXj42XBUDOa-9-YVrfz5nNDObPDoaPcyYL032NVmA_IUC9xw2FUV63KB_rk-SupPbjK2jiQ6viK51NIBf4l_d1hxt8fIEsBdGQUNiyCsJhhevt_flIccIJ5Y1mTaF07FfgbTvA",
			wantErr:     false,
			wantTimeout: true,
			timeout:     2,
		},
		{
			name: "005",
			fields: fields{
				opt: JwsHelperOpt{
					Audience:       "xyt",
					Issuer:         "gw123",
					Timeout:        1,
					PublicKeyPath:  "./test/1_rsa_public_key.pem",
					PrivateKeyPath: "./test/1_rsa_private_key.pem",
					HashIdsSalt:    "",
				},
			},
			args: args{
				userId:  81983918293,
				payload: "hello world",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jws, err := NewJwsHelper(tt.fields.opt)
			if err != nil {
				t.Fatal(err)
			}
			got, err := jws.CreateToken(tt.args.userId, tt.args.payload)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if got != tt.want {
			//	t.Errorf("CreateToken() got = %v \n want %v", got, tt.want)
			//}
			if tt.wantTimeout {
				time.Sleep(time.Duration(tt.timeout) * time.Second)
			}
			userId, payload, err := jws.ParseToken(got)
			if tt.wantTimeout {
				if err == nil {
					t.Errorf("ParseToken() wantTimeoutError, userId :%d ,payload :%s", userId, payload)
				}
				t.Log("debug ==>", err)
			} else {
				if err != nil || userId != tt.args.userId || payload != tt.args.payload {
					t.Errorf("ParseToken() error = %v, userId :%d ,payload :%s", err, userId, payload)
				}
			}
			t.Log(userId, payload)
		})
	}
}
