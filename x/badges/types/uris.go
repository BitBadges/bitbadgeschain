package types

//TODO: test this actually works correctly
func GetSubassetUriFromUriObject(uriObject UriObject) (string, error) {
	uri, err := GetUriFromUriObject(uriObject)
	if err != nil {
		return "", err
	}

	subassetUri := uri[:uriObject.IdxRangeToRemove.Start] + uri[uriObject.IdxRangeToRemove.End:]
	subassetUri = subassetUri[:uriObject.InsertSubassetBytesIdx] + string(uriObject.BytesToInsert) + subassetUri[uriObject.InsertSubassetBytesIdx:]
	subassetUri = subassetUri[:uriObject.InsertIdIdx] + "0" + subassetUri[uriObject.InsertIdIdx:]

	return subassetUri, nil
}


func GetUriFromUriObject(uriObject UriObject) (string, error) {
	uri := ""
	if uriObject.Scheme != 0 {
		if uriObject.Scheme == 1 {
			uri += "http://"
		} else if uriObject.Scheme == 2 {
			uri += "https://"
		} else if uriObject.Scheme == 3 {
			uri += "ipfs://"
		} else {
			return "", ErrInvalidUriScheme
		}
	}

	uri += string(uriObject.Uri)
	return uri, nil
}