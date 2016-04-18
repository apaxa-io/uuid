package uuid

import "bytes"

// Returns JSON valid representation of a uuid
func (u UUID) MarshalJSON() ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteString("\"")
	buffer.WriteString(u.String())
	buffer.WriteString("\"")
	return []byte(buffer.String()), nil
}

/*
//TODO compare perfomance
func (u UUID) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.String())
}
*/

// Parses a JSON representation of a uuid
func (u *UUID) UnmarshalJSON(bytes []byte) error {
	p, err := ParseHex(string(bytes))
	if err != nil {
		return err
	}
	copy(u[:], p[:])
	return nil
}

/*
//TODO compare perfomance
func (u *UUID) UnmarshalJSON(data []byte) error {
	var uuid_as_string string

	if u == nil {
		u = new(UUID)
	}

	if err := json.Unmarshal(data, &uuid_as_string); err != nil {
		return err
	}

	parsed_uuid, err := ParseHex(uuid_as_string)
	if err != nil {
		return err
	}

	*u = *parsed_uuid

	return nil
}
*/
