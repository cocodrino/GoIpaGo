package main

import "testing"

func TestPronounce(t *testing.T){
	if texto:=Pronounce(Ipa,"hello how are you");texto != "hʌl̩oʊ haʊ ɑɹ ju"{
		t.Errorf("Ipa hello how are you: expected hʌl̩oʊ haʊ ɑɹ ju received %s",texto)
	}

	if texto:=Pronounce(Simplified,"hello how are you");texto != "Jʌl̩Ou JAu Ar yU"{
		t.Errorf("Simplified hello how are you: expected Jʌl̩Ou JAu Ar yU received %s",texto)
	}

	if texto:=Pronounce(Simplified,"this text include , and ! characters");texto != "Dis tEkst in̩kl̩Ud , ʌn̩d ! kariktorz"{
		t.Errorf("Simplified with , and !: expected\n Dis tEkst in̩kl̩Ud , ʌn̩d ! kariktorz\nreceiv: %s",texto)
	}

	if texto:=Pronounce(Simplified,"this is not a word bazinga !");texto != "Dis iz n̩At ʌ word bazinga !"{
		t.Errorf("Simplified with , and !: expected\n Dis iz n̩At ʌ word bazinga ! \nreceived\n %s",texto)
	}

}
