package comparator

import (
	"math"
)

type FFT struct{
	 m int
	 n int
	 cos []float64
	 sin []float64
	 window	[]float64
}


func GetInstance(n int)*FFT{
	m:=uint(math.Log(float64(n))/math.Log(float64(2)));
	
	if n!=(1<<m){
		panic("FFT length must be power of 2")
	}
	
	cos := make([]float64,n/2)
	sin := make([]float64,n/2)
	
	for i :=0;i<n/2;i++{
		cos[i] = math.Cos(-2*math.Pi*float64(i)/float64(n));
		sin[i] = math.Sin(-2*math.Pi*float64(i)/float64(n));
	}
	window:=makeWindow(n)
	return &FFT{
		m:	int(m),
		n:	n,
		cos: cos,
		sin: sin,
		window: window,
	}
}

func makeWindow(n int)[]float64{
	window:=make([]float64,n)
	for i:=0;i<n;i++{
		window[i]=0.42 - 0.5*math.Cos(2*math.Pi*float64(i)/float64(n-1)) +
			0.08*math.Cos(4*math.Pi*float64(i)/float64(n-1))
	}
	return window
}

func (fft *FFT) GetWindow()[]float64{
	return fft.window
}


  /***************************************************************
  * fft.go
  *   fft: in-place radix-2 DIT DFT of a complex input 
  * 
  *   input: 
  * n: length of FFT: must be a power of two 
  * m: n = 2**m 
  *   input/output 
  * x: double array of length n with real part of data 
  * y: double array of length n with imag part of data 
  * 
  ****************************************************************/
  
  func (fft *FFT) Fft(x,y []float64){
  	var i int
  	var j int
  	var k int
  	var n1 int
  	var n2 int
  	var a int
  	
  	var c float64
  	var s float64
//  	var e float64
  	var t1 float64
  	var t2 float64
  	
  	j=0
  	n2=fft.n/2;
  	for i=1;i<fft.n-1;i++{
  		n1=n2;
  		for ;j>=n1;{
  			j=j-n1
  			n1=n1/2
  		}
  		j=j+n1
  		if i<j{
  			t1=x[i]
  			x[i]=x[j]
  			x[j]=t1
  			t1=y[i]
  			y[i]=y[j]
  			y[j]=t1
  		}
  	}
  	//FFT
  	n1=0
  	n2=1
  	for i=0;i<fft.m;i++{
  		n1=n2
  		n2=n2+n2
  		a=0
  		
  		for j=0;j<n1;j++{
  			c=math.Cos(float64(a))
  			s=math.Sin(float64(a))
  			a +=1<<uint(fft.m-i-1)
  			for k=j;k<fft.n;k=k+n2{
  				t1=c*x[k+n1] - s*y[k+n1]
  				t2=s*x[k+n1] + c*y[k+n1]
  				x[k+n1]=x[k] - t1
  				y[k+n1]=y[k] - t2
  				x[k]=x[k]+t1
  				y[k]=y[k]+t2
  			}
  		}
  	}
  	
  }