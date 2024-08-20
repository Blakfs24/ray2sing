package ray2sing

import (
	T "github.com/sagernet/sing-box/option"
)

type outboundMap map[string]interface{}

// type outboundMap map[string]Any

func VlessXray(vlessURL string) (*T.Outbound, error) {
	u, err := ParseUrl(vlessURL, 443)
	if err != nil {
		return nil, err
	}
	decoded := u.Params
	// fmt.Printf("Port %v deco=%v", port, decoded)
	streamSettings, err := getStreamSettingsXray(decoded)
	if err != nil {
		return nil, err
	}

	// packetEncoding := decoded["packetencoding"]
	// if packetEncoding==""{
	// 	packetEncoding="xudp"
	// }

	return &T.Outbound{
		Tag:  u.Name,
		Type: "xray",
		XrayOptions: T.XrayOutboundOptions{
			// DialerOptions: getDialerOptions(decoded),
			Fragment: getXrayFragmentOptions(decoded),
			XrayOutboundJson: &map[string]any{

				"protocol": "vless",
				"settings": map[string]any{
					"vnext": []any{
						map[string]any{
							"address": u.Hostname,
							"port":    u.Port,
							"users": []any{
								map[string]string{
									"id":         u.Username, // Change to your UUID.
									"encryption": "none",
									"flow":       decoded["flow"],
								},
							},
						},
					},
				},
				"tag":            u.Name,
				"streamSettings": streamSettings,
				"mux":            getMuxOptionsXray(decoded),
			},
		},
	}, nil
}

func Singbox2XrayVless(out map[string]any) (*T.Outbound, error) {
	streamSettings, err := getStreamSettingsXrayFromSingbox(out)
	if err != nil {
		return nil, err
	}
	mux := map[string]any{
		"enabled":     true,
		"concurrency": 0,
		// "xudpConcurrency": 16,
		// "xudpProxyUDP443": "reject"
	}
	config := T.Outbound{
		Tag:  out["tag"].(string),
		Type: "xray",
		XrayOptions: T.XrayOutboundOptions{
			// Fragment: getXrayFragmentOptions(decoded),
			XrayOutboundJson: &map[string]any{
				"protocol": "vless",
				"settings": map[string]any{
					"vnext": []any{
						map[string]any{
							"address": out["server"].(string),
							"port":    out["server_port"].(float64),
							"users": []any{
								map[string]string{
									"id":         out["uuid"].(string), // Change to your UUID.
									"encryption": "none",
									"flow":       "",
								},
							},
						},
					},
				},
				"tag":            out["tag"].(string),
				"streamSettings": streamSettings,
				"mux":            mux,
			},
		},
	}
	return &config, nil
}
