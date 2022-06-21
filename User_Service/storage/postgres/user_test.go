package postgres

import (
	"reflect"
	"testing"

	pb "github.com/mahmud3253/Project/User_Service/genproto"
)
	
func TestUserRepo_Crete(t *testing.T) {
	tests := []struct {
		name    string
		input   *pb.User
		want    *pb.User
		wantErr bool
	}{
		{
			name: "success case",
			input: &pb.User{
				FirstName: "Mahmud",
				LastName:  "Boriyev",
				Posts:     nil,
			},
			want: &pb.User{
				FirstName: "Mahmud",
				LastName:  "Boriyev",
				Posts:     nil,
			},
			wantErr: false,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := repo.CreateUser(tc.input)
			if err != nil {
				t.Fatalf("%s: expected: %v, got: %v", tc.name, tc.wantErr, err)
			}
			got.Id = ""
			//got.Posts = nil
			if !reflect.DeepEqual(tc.want, got) {
				t.Fatalf("%s: expected: %v,got: %v", tc.name, tc.want, got)
			}
		})
	}
}
			