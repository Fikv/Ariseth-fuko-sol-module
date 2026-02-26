package domain

type link struct {
	linkType string;
	label string; 
	url string;
}


type token struct{
	url string;
	chainId string;
	tokenAdress string;
	icon string; 
	header string; 
	description string; 
	links []link;
}