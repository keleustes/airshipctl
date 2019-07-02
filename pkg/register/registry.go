package register

import ()

func RegisterCoreKinds() {
	if err := NewArmadaCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewKfDefCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewNetworkingCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	// if err := NewOpenShiftCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
	//	panic(err)
	// }
	if err := NewRolloutCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewWorkflowCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewClusterApiCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewBootstrapProviderKubeadmCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	if err := NewBareMetalProviderCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
		panic(err)
	}
	// if err := NewOpenStackProviderCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
	//	panic(err)
	// }
	// if err := NewAWSProviderCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
	//	panic(err)
	// }
	// if err := NewVSphereProviderCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
	//	panic(err)
	// }
	// if err := NewDockerProviderCRDRegisterPlugin().Config(nil, []byte{}); err != nil {
	//	panic(err)
	// }
}
