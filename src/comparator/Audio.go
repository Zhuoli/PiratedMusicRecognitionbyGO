package comparator

import (
	"strings"
	"bufio"
	"errors"
	"os"
	"bytes"
	"encoding/binary"
	 "math"
)


type Audio struct{
	fileName string
	sampleRate float64
	waveFormat	string
	bitesPerSecond int
	
	headerArray []byte
	dataArray []byte	

	leftChannelSamples []int16
	timeZoneData []float64
	frequenciesData []float64
}

func NewAudio(path string) *Audio{
	var name string
	var dataArray []byte
	var headerArray []byte
	var sampleRate float64
	var waveFormat	string
	var bitesPerSecond int
	
	strs :=strings.Split(path, "/")
	name =strs[len(strs)-1]
	
	f,err :=os.Open(path)
	if  err!=nil{
		panic(err)
	}
	defer f.Close()	
	reader:=bufio.NewReader(f)
	size :=reader.Buffered()
	headerArray =make([]byte,44)
	_,err=reader.Read(headerArray)
	if err!=nil{
		panic(err)
	}
	sampleRate,waveFormat,bitesPerSecond = checkHeaderFormat(headerArray)
	dataArray =make([]byte,size-44)
	_,err =reader.Read(dataArray)
	var fileLeftChannel []byte = extractLeftChannels(dataArray)
	var timeZoneData []float64 = convert2doubles(fileLeftChannel)
	var leftChannelSamples []int16 = convert2short(fileLeftChannel)
	var fileImg []float64 = applyFFT(timeZoneData)
	var frequenciesData []float64=convert2frequencies(fileImg,timeZoneData)
	
	return &Audio{
		fileName: name,
		sampleRate: sampleRate,
		waveFormat: waveFormat,
		bitesPerSecond: bitesPerSecond,
		
		headerArray:	headerArray,
		dataArray:	    dataArray,
		leftChannelSamples:	leftChannelSamples,
		timeZoneData:	timeZoneData,
		frequenciesData: frequenciesData,
	}
	
}

func checkHeaderFormat(bytes []byte)(float64,string,int){
	if bytes[8]!=87 || bytes[9]!=65 || bytes[10] !=86 || bytes[11]!=69{
		err :=errors.New("invalid format: not wav")
		panic(err)
	}
	if bytes[20]!=1 || bytes[21]!=0{
		err:=errors.New("invalid format: not PCM")
		panic(err)
	}
	if(bytes[22]!=2 || bytes[23]!=0){
		err:=errors.New("invalid format: channels not STEREO")
		panic(err)
	}
	if(bytes[24]!=68||bytes[25]!=172||bytes[26]!=0||bytes[27]!=0){
		err:=errors.New("invalid format: not 44100 sample rate")
		panic(err)
	}
	if(bytes[34]!=16||bytes[35]!=0){
		err:=errors.New("invalid format: Not 16 bites per sample")
		panic(err)
	}
	return 44.1,"wav",16;
}

func extractLeftChannels(dataArray []byte)[]byte{
	fileLeftChannel := make([]byte,len(dataArray)/2)
	for i:=0;i<len(fileLeftChannel)/2;i++{
		fileLeftChannel[i*2]=dataArray[i*4]
		fileLeftChannel[i*2+1]=dataArray[i*4+1];
	}
	return fileLeftChannel;
}
func convert2doubles(fileLeftChannel []byte)[]float64{
	fileDouble := make([]float64,len(fileLeftChannel)/2)
	buffer :=bytes.NewBuffer(fileLeftChannel)
	for i:=0;;i++ {
		twoBytes:=buffer.Next(2);
		if len(twoBytes)!=2{
			break;	
		}
		t,n:=binary.Varint(twoBytes)
		if n<=0{
			err:=errors.New("short byte error")
			panic(err)
		}
		fileDouble[i]=float64(t)/float64(32768.0)	
	}
	return fileDouble;
	
}
func applyFFT(timeZoneData []float64)[]float64{
	return nil;
}
func convert2frequencies(fileImg,fileDouble []float64)[]float64{
	frequenciesData :=make([]float64,len(fileDouble)/2)
	for j:=0;j<len(frequenciesData);j++{
		reql:=fileDouble[j]
		img:=fileImg[j]
		freq:=math.Sqrt(reql*reql+img*img)
		frequenciesData[j]=freq
	}
	return frequenciesData
	
}
func convert2short(fileLeftChannel []byte)[]int16{
	shortArray := make([]int16,len(fileLeftChannel)/2)
	buffer :=bytes.NewBuffer(fileLeftChannel)
	for i:=0;;i++ {
		twoBytes:=buffer.Next(2);
		if len(twoBytes)!=2{
			break;	
		}
		t,n:=binary.Varint(twoBytes)
		if n<=0{
			err:=errors.New("short byte error")
			panic(err)
		}
		shortArray[i]=int16(t)
	}
	return shortArray;
}