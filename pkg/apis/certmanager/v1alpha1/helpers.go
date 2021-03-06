package v1alpha1

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (i *IssuerStatus) ACMEStatus() *ACMEIssuerStatus {
	if i.ACME == nil {
		i.ACME = &ACMEIssuerStatus{}
	}
	return i.ACME
}

func (a *ACMEIssuerDNS01Config) Provider(name string) (*ACMEIssuerDNS01Provider, error) {
	for _, p := range a.Providers {
		if p.Name == name {
			return &(*&p), nil
		}
	}
	return nil, fmt.Errorf("provider '%s' not found", name)
}

func (a *ACMECertificateConfig) ConfigForDomain(domain string) ACMECertificateDomainConfig {
	for _, cfg := range a.Config {
		for _, d := range cfg.Domains {
			if d == domain {
				return cfg
			}
		}
	}
	return ACMECertificateDomainConfig{}
}

func (c *CertificateStatus) ACMEStatus() *CertificateACMEStatus {
	if c.ACME == nil {
		c.ACME = &CertificateACMEStatus{}
	}
	return c.ACME
}

func (c *CertificateACMEStatus) SaveAuthorization(a ACMEDomainAuthorization) {
	for i, auth := range c.Authorizations {
		if auth.Domain == a.Domain {
			c.Authorizations[i] = a
			return
		}
	}
	c.Authorizations = append(c.Authorizations, a)
}

func (iss *Issuer) HasCondition(condition IssuerCondition) bool {
	if len(iss.Status.Conditions) == 0 {
		return false
	}
	for _, cond := range iss.Status.Conditions {
		if condition.Type == cond.Type && condition.Status == cond.Status {
			return true
		}
	}
	return false
}

func (iss *Issuer) UpdateStatusCondition(conditionType IssuerConditionType, status ConditionStatus, reason, message string) {
	newCondition := IssuerCondition{
		Type:    conditionType,
		Status:  status,
		Reason:  reason,
		Message: message,
	}

	t := time.Now()

	if len(iss.Status.Conditions) == 0 {
		glog.Infof("Setting lastTransitionTime for Issuer %q condition %q to %v", iss.Name, conditionType, t)
		newCondition.LastTransitionTime = metav1.NewTime(t)
		iss.Status.Conditions = []IssuerCondition{newCondition}
	} else {
		for i, cond := range iss.Status.Conditions {
			if cond.Type == conditionType {
				if cond.Status != newCondition.Status {
					glog.Infof("Found status change for Issuer %q condition %q: %q -> %q; setting lastTransitionTime to %v", iss.Name, conditionType, cond.Status, status, t)
					newCondition.LastTransitionTime = metav1.NewTime(t)
				} else {
					newCondition.LastTransitionTime = cond.LastTransitionTime
				}

				iss.Status.Conditions[i] = newCondition
				break
			}
		}
	}
}

func (crt *Certificate) HasCondition(condition CertificateCondition) bool {
	if len(crt.Status.Conditions) == 0 {
		return false
	}
	for _, cond := range crt.Status.Conditions {
		if condition.Type == cond.Type && condition.Status == cond.Status {
			return true
		}
	}
	return false
}

func (crt *Certificate) UpdateStatusCondition(conditionType CertificateConditionType, status ConditionStatus, reason, message string) {
	newCondition := CertificateCondition{
		Type:    conditionType,
		Status:  status,
		Reason:  reason,
		Message: message,
	}

	t := time.Now()

	if len(crt.Status.Conditions) == 0 {
		glog.Infof("Setting lastTransitionTime for Certificate %q condition %q to %v", crt.Name, conditionType, t)
		newCondition.LastTransitionTime = metav1.NewTime(t)
		crt.Status.Conditions = []CertificateCondition{newCondition}
	} else {
		for i, cond := range crt.Status.Conditions {
			if cond.Type == conditionType {
				if cond.Status != newCondition.Status {
					glog.Infof("Found status change for Certificate %q condition %q: %q -> %q; setting lastTransitionTime to %v", crt.Name, conditionType, cond.Status, status, t)
					newCondition.LastTransitionTime = metav1.NewTime(t)
				} else {
					newCondition.LastTransitionTime = cond.LastTransitionTime
				}

				crt.Status.Conditions[i] = newCondition
				break
			}
		}
	}
}
