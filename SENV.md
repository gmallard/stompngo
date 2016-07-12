# senv - A stompngo helper (sub) package

The intent of this package is to assist the clients of the
_stompngo_ package when developing and using _stompngo_ client application
code.

This assistance is _primarily_ implemented in the form of environment variables:

* export=STOMP_var=value (Unix systems)
* set STOMP_var=value (Windows systems)

Using these environment variables can help avoid or eliminate 'hard coding'
values in the client code.

Environment variables and related subjects are discussed in the following sections:

* [Supported Environment Variables](#senv)
* [Supported Helper Function](#shf)
* [Example Code Fragments](#ecf)
* [Complete Connect Header Fragment](#cchf)
<br />

---

## <a name="sev"></a>Supported Environment Variables

The following table shows currently supported environment variables.

<table border="1" style="width:80%;border: 1px solid black;">
<tr>
<th style="width:20%;border: 1px solid black;padding-left: 10px;" >
Environment Variable Name
</th>
<th style="width:60%border: 1px solid black;padding-left: 10px;" >
Usage
</th>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_DEST
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A destination to be used by the client.<br />
Default: /queue/sng.sample.stomp.destination
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_LOGIN
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The login to be used by the client in the CONNECT frame.<br />
Default: guest
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_HEARTBEATS
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
For protocol 1.1+, the heart-beat value to be used by the client in the CONNECT frame.<br />
Default: 0,0
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_HOST
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The broker host to connect to.<br />
Default: localhost
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_NMSGS
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
A default nummber of messages to receive or send.  Useful for some clients.<br />
Default: 1
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_PASSCODE
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The passcode to be used by the client in the CONNECT frame.<br />
Default: guest
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_PERSISTENT
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
May control use of the persistent header in SEND frames.<br />
Default: no persistent header to be used.<br />
Example:<br />
STOMP_PERSISTENT=anyvalue
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_PORT
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
The broker port to connect to.<br />
Default: 61613
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_PROTOCOL
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
For protocol 1.1+, the accept-version value to be used by the client in the CONNECT frame.<br />
Default: 1.2<br />
Multiple versions may be used per the STOMP specifications, e.g.:<br />
STOMP_PROTOCOL="1.0,1.1,1.2"
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_SUBCHANCAP
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
Used to possibly override the default capacity of _stompngo_ subscription channels.<br />
Default: 1 (the same as the _stompngo_ default.)
</td>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
STOMP_VHOST
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
For protocol 1.1+, the host value to be used by the client in the CONNECT frame.<br />
Default: If not specified the default is STOMP_HOST (i.e. localhost).
</td>
</tr>
</table>

<br />

---

## <a name="shf"></a>Supported Helper Function

There is currently one helper function also provided by the _senv_ package.

<table border="1" style="width:80%;border: 1px solid black;">
<tr>
<th style="width:20%;border: 1px solid black;padding-left: 10px;" >
Function Name
</th>
<th style="width:60%border: 1px solid black;padding-left: 10px;" >
Usage
</th>
</tr>

<tr>
<td style="border: 1px solid black;padding-left: 10px;" >
HostAndPort()
</td>
<td style="border: 1px solid black;padding-left: 10px;" >
This function returns two values, the STOMP_HOST and STOMP_PORT values.<br />
Example:<br />
h, p := senv.HostAndPort()
</td>
</tr>

</table>

<br />

---

## <a name="ecf"></a>Example Code Fragments

Example code fragments follow.

### STOMP_DEST Code Fragment

        sbh := ..... // Subscribe header, type stompngo.Headers
        if senv.Dest() != "" {
            sbh = sbh.Add("destination", senv.Dest())
        }

### STOMP_HEARTBEATS Code Fragment

        ch := ..... // Connect headers, type stompngo.Headers
        if senv.Heartbeats() != "" {
            ch = ch.Add("heart-beat", senv.Heartbeats())
        }

### STOMP_HOST Code Fragment

        ch := ..... // Connect headers, type stompngo.Headers
        if senv.Host() != "" {
            ch = ch.Add("host", senv.Host())
        }

### STOMP_LOGIN Code Fragment

        ch := ..... // Connect headers, type stompngo.Headers
        if senv.Login() != "" {
            ch = ch.Add("login", senv.Login())
        }

### STOMP_NMSGS Code Fragment

    msg_count := senv.Nmsgs() // Default is 1

### STOMP_PASSCODE Code Fragment

        ch := ..... // Connect headers, type stompngo.Headers
        if senv.Passcode() != "" {
            ch = ch.Add("passcode", senv.Passcode())
        }

### STOMP_PERSISTENT Code Fragment

        sh := ..... // Send headers, type stompngo.Headers
        if senv.Persistent() != "" {
            ch = ch.Add("persistent", "true") // Brokers might need 'true' here
        }

### STOMP_PORT Code Fragment

        ch := ..... // Connect headers, type stompngo.Headers
        if senv.Port() != "" {
            ch = ch.Add("port", senv.Port())
        }

### STOMP_PROTOCOL Code Fragment

        ch := ..... // Connect headers, type stompngo.Headers
        if senv.Protocol() != "" {
            ch = ch.Add("accept-version", senv.Protocol())
        }

### STOMP_VHOST Code Fragment

        ch := ..... // Connect headers, type stompngo.Headers
        if senv.Vhost() != "" {
            ch = ch.Add("host", senv.Vhost())
        }
<br />

---

## <a name="cchf"></a>Complete Connect Header Fragment

Obtaining a full set of headers to use for a _stompngo.Connect_ might
look like this:

        func ConnectHeaders() stompngo.Headers {
          h := stompngo.Headers{}
          l := senv.Login()
          if l != "" {
              h = h.Add("login", l)
          }
          pc := senv.Passcode()
          if pc != "" {
              h = h.Add("passcode", pc)
          }
          //
          p := senv.Protocol()
          if p != stompngo.SPL_10 { // 1.1 and 1.2
              h = h.Add("accept-version", p).Add("host", senv.Vhost())
          }
          //
          hb := senv.Heartbeats()
          if hb != "" {
              h = h.Add("heart-beat", hb)
          }
          return h
        }

