// todo we should move this to oses as platform specific
package ssh

type RemoteInfo struct {
	OS                string `json:"os" yaml:"os"`
	Hostname          string `json:"hostname" yaml:"hostname"`
	KernelVersion     string `json:"kernel_version" yaml:"kernel_version"`
	KernelRelease     string `json:"kernel_release" yaml:"kernel_release"`
	Architecture      string `json:"arch" yaml:"arch"`
	HetznerRescueMode string `json:"hetzner_rescue_mode" yaml:"hetzner_rescue_mode"`
}

func (s *Client) RemoteInfo() (*RemoteInfo, error) {
	if s.info != nil {
		return s.info, nil
	}
	if _, err := s.Connect(); err != nil {
		return nil, err
	}

	info := RemoteInfo{
		OS:            "uname",
		Hostname:      "uname -n",
		KernelVersion: "uname -v",
		KernelRelease: "uname -r",
		Architecture:  "uname -m",
		HetznerRescueMode: `{ 
			grep -q "Hetzner Rescue System" /etc/motd &>/dev/null && 
			test -f /root/.oldroot/nfs/install/installimage &&
			echo yes || 
			echo
		}`,
	}
	if err := s.OutputPtrs([]*string{
		&info.OS,
		&info.Hostname,
		&info.KernelVersion,
		&info.KernelRelease,
		&info.Architecture,
		&info.HetznerRescueMode,
	}); err != nil {
		return nil, err
	}

	s.info = &info
	return &info, nil
}
