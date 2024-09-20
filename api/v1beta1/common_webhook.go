package v1beta1

const (
	// ErrPrivilegedModeRequired
	ErrPrivilegedModeRequired = "%s.Spec.Privileged is requied in order to successfully " +
	                            "execute tests with the provided configuration."
)

const (
	// WarnPrivilegedModeOn
	WarnPrivilegedModeOn = "%s.Spec.Privileged is set to true. This means that test pods " +
												 "are spawned with allowPrivilegedEscalation: true and default " +
												 "capabilities on top of those required by the test operator " +
												 "(NET_ADMIN, NET_RAW)."

	// WarnPrivilegedModeOff
	WarnPrivilegedModeOff = "%[1]s.Spec.Privileged is set to false. Note, that a certain " +
	                        "set of tests might fail, as this configuration may be " +
													"required for the tests to run successfully. Before enabling" +
													"this parameter, consult documentation of the %[1]s CR."
)