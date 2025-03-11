package main

import (
	"fmt"

	"github.com/telexintegrations/support-ai/api"
	"github.com/telexintegrations/support-ai/format"
)

func main() {
	config, _ := api.LoadEnvConfig()

	chunkThis := "Cybersecurity is a crucial field in the modern digital era, where data breaches and cyber threats have become increasingly sophisticated. Organizations and individuals rely on various security measures to safeguard sensitive information from unauthorized access, theft, and manipulation. One of the most fundamental aspects of cybersecurity is encryption, a technique that ensures data remains confidential even if intercepted by malicious actors. Encryption works by converting plain text into an unreadable format using mathematical algorithms. Only authorized parties possessing the correct decryption key can revert the data to its original form. There are two primary types of encryption: symmetric encryption, where the same key is used for both encryption and decryption, and asymmetric encryption, which involves a pair of public and private keys. Popular encryption standards include AES (Advanced Encryption Standard), RSA (Rivest-Shamir-Adleman), and ECC (Elliptic Curve Cryptography). The importance of encryption extends beyond just protecting personal or corporate data. It is widely used in securing financial transactions, ensuring the integrity of communication channels, and safeguarding cloud storage. For instance, online banking and e-commerce platforms rely on encryption protocols such as TLS (Transport Layer Security) to secure customer transactions and prevent unauthorized access to sensitive financial data. Cybercriminals continuously develop new tactics to bypass security measures, making it necessary for organizations to adopt a multi-layered security approach. This includes encryption, multi-factor authentication (MFA), intrusion detection systems (IDS), and endpoint security solutions. AI-driven cybersecurity tools further enhance protection by analyzing large volumes of data to identify suspicious activities and potential threats in real-time. Another emerging trend in cybersecurity is the adoption of Zero Trust Architecture (ZTA), which operates on the principle that no user or device should be trusted by default. This approach enforces strict identity verification, continuous monitoring, and least-privilege access policies to minimize security risks. Organizations implementing Zero Trust benefit from enhanced visibility and control over their digital assets.With the growing reliance on cloud computing, data security has become more complex. Cloud service providers implement robust security measures, including encryption, access control, and automated threat detection. However, businesses must also take proactive steps to secure their cloud environments by ensuring data is encrypted at rest and in transit, regularly auditing access controls, and implementing identity and access management (IAM) policies. In conclusion, encryption remains a cornerstone of modern cybersecurity strategies. As cyber threats continue to evolve, organizations and individuals must stay vigilant and adopt best security practices to protect sensitive data. Implementing encryption alongside advanced security measures, AI-driven threat detection, and Zero Trust principles will help mitigate risks and safeguard the digital world against emerging threats."
	chunkedResp := format.ChunkTextByParagraph(chunkThis, 30)
	
	for i, chunk := range chunkedResp {
		fmt.Printf("Chunk %d:\n%s\n\n--Here--\n\n", i+1, chunk)
	}
	
	server := api.NewServer(&config)
	server.StartServer(":8080")
}
