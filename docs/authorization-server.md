# All about Authorization Server

Authorization server adalah hal yang paling komplex, kenapa? karena client service dan protected service akan bergantung ke authorization server.

## Endpoint

- /authorize
- /approve
- /token

### Authorize Endpoint

Endpoint authorize yang akan dipanggil pertama kali oleh client. Tugasnya untuk mengecek apakah client tersebut terdaftar atau tidak diserver authorization. Jika benar ada, maka akan merender halaman approval. Dihalaman approval ada nilai tersembunyi yang digunakan untuk pengamanan permintaan, diletakan dengan nama request ID, teknik ini dikenal dengan CSRF. Request ID hanya berupa random string yang disimpan diserver yang nantinya akan dicocokan.

#### Apa saja yang perlu diperiksa oleh endpoint Authorize

1. ResponseType = [https://tools.ietf.org/html/rfc6749#section-3.1.1](https://tools.ietf.org/html/rfc6749#section-3.1.1)
2. ClientID     = client id
3. Redirect URI = Endpoint untuk menerima balikan dari authorization server
4. State        = Kode unik yang digunakan client untuk verifikasi

#### Apa saja yang perlu dicatat oleh endpoint Authorize

  1. Client ID
  2. Clinet Secret
  3. RedirectURI
  4. Requests          = Requests dari client ke authorize dengan id requestID

Kenapa kita perlu mencatat request dari client, karena nantinya kita membutuhkan statenya untuk dikembalikan ke redirect uri mereka.

Kurang lebih seperti ini:

```js
// disimpan sebagai array
requests["eycJPJHY"] = {
    "client_id": "oauth-client-1",
    "redirect_uri": "http://localhost:8081/callback",
    "response_type":"code",
    "state":"XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDa"
}
```

#### Apa yang dihasilkan ketika akses endpoint Authorize

1. Halaman Approval
2. Request ID

#### Dimana menyimpan Request ID

Bisa ditempat penyimpanan data apapun. Sifat request ID adalah sementara, Setelah dicocokan akan langsung dihapus.

### Approve Endpoint

Approve endpoint digunakan untuk melakukan pemerikasan penyetujuan dari resource owner. Approve hanya memberikan 2 pilihan, Setuju atau Tolak. Jika Setuju maka akan diredirect kehalaman `RedirectURI` client. Ketika redirect akan disertakan code dan state yang diambil dari request yang disimpan pada enpoint Authorize sebelumnya.

#### Apa saja yang perlu diperiksa oleh endpoint Approve

1. Cari requests berdasarkan request id
2. Apakah pengguna menyetujui atau menolak, bisa dilihat dari membaca form htmlnya
3. Apakah response typenya `code`

#### Apa saja yang perlu dicatat oleh endpoint Approve

1. code, ini mirip seperti request id

Kenapa kita perlu mencatat code yang berisi query, karena nantinya akan digunakan ketika pengambilan token

Kurang lebih seperti ini yang disimpan:

```js
// disimpan sebagai array
code["eycJPJHY"] = {
    "requests": {
        "client_id": "oauth-client-1",
        "redirect_uri": "http://localhost:8081/callback",
        "response_type":"code",
        "state":"XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDa"
    }
}
```

#### Apa yang dihasilkan ketika akses endpoint Approve

Redirect ke Redirecy URI `http://localhost:8081/callback?code=NufNjJhh&state=XVlBzgbaiCMRAjWwhTHctcuAxhxKQFDa`

- Code  = yang akan digunakan untuk mengambil token
- State = yang akan digunakan untuk pemeriksaan oleh client

### Token Endpoint

Enpoint ini bertugas