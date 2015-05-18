package httpauth

import (
    "log"
//	"errors"
	"github.com/gocql/gocql"
//	"github.com/elvtechnology/gocqltable"
//	"github.com/elvtechnology/gocqltable/recipes"
	//"gopkg.in/mgo.v2"
	//"gopkg.in/mgo.v2/bson"
)

// MongodbAuthBackend stores database connection information.
//type MongodbAuthBackend struct {
//	mongoURL string
//	database string
//	session  *mgo.Session
//}
type CassandraAuthBackend struct {
	cassandraURLs []string
	keyspace string
	consistency gocql.Consistency
	cluster *gocql.ClusterConfig//gocql.Cluster
	session  *gocql.Session
}
/*
func (b MongodbAuthBackend) connect() *mgo.Collection {
	session := b.session.Copy()
	return session.DB(b.database).C("goauth")
}

func mkmgoerror(msg string) error {
	return errors.New("mongobackend: " + msg)
}*/

// NewMongodbBackend initializes a new backend.
// Be sure to call Close() on this to clean up the mongodb connection.
// Example:
//     backend = httpauth.MongodbAuthBackend("mongodb://127.0.0.1/", "auth")
//     defer backend.Close()
func NewCassandraBackend(cassandraURLs []string, k string, c gocql.Consistency) (b CassandraAuthBackend, err error) {
	b.cassandraURLs	= cassandraURLs
	b.keyspace		= k
	b.consistency	= c
    b.cluster	   = gocql.NewCluster(cassandraURLs...)
	b.session, err = b.cluster.CreateSession()
	if err != nil {
		log.Fatalln("Unable to open up a session with Cassandra (err="+err.Error() + ")")
	}

    return b, err;
	/*gocqltable.SetDefaultSessions(b.session)

	keyspace := gocqltable.NewKeyspace(k)

	err = keyspace.Create(map[string]interface{}{
		"class":	"SimpleStrategy",
		"replication_factor": 1,
	}, true)
	if err != nil {
		log.Fatalln(err)
	}

	type User struct {
		Username	string
		Email		string
		Hash		string
		Role		string
		Created		time.Time
	}*/

	/*goauthTable := struct {
		recipes.CRUD
	}{
		recipes.CRUD{
			b.keyspace.NewTable(
				"goauth",
				[]string{"username"},
				nil,
				User{},
			),
		},
	}*/

	//err = goauthTable.Create()
	/*if err != nil {
		log.Fatalln(err)
	}*/
}

func (b CassandraAuthBackend) User(username string) (user UserData, e error) {
    session := b.session
    if err := session.Query(`SELECT username, email, hash, role FROM users WHERE username = ? LIMIT 1`,
    username).Consistency(gocql.One).Scan(&user.Username, &user.Email, &user.Hash, &user.Role); err != nil {
        log.Fatal(err)
        return user, err
    }
    user.Username = username
    return user, nil
}

// Users returns a slice of all users.
func (b CassandraAuthBackend) Users() (us []UserData, e error) {
    var (
        username, email, role string
        hash                  []byte
    )
    session := b.session
    iter := session.Query(`SELECT username, email, hash, role FROM users`).Iter()
    next := iter.Scan(&username, &email, &hash, &role)
    for next {
        us = append(us, UserData{username, email, hash, role})
    }
    return us, nil
}

// SaveUser adds a new user, replacing one with the same username.
func (b CassandraAuthBackend) SaveUser(user UserData) (err error) {
    session := b.session
    if _, err := b.User(user.Username); err == nil {
        if err = session.Query(`INSERT INTO users (username, email, hash, role) VALUES (?, ?, ?, ?)`,
        user.Username, user.Email, user.Hash, user.Role).Exec(); err != nil {
            log.Fatal(err)
            return err
        }
    } else {
        if err = session.Query(`UPDATE users SET email=? hash=? role=?) VALUES (?, ?, ?) WHERE username=?`,
        user.Email, user.Hash, user.Role, user.Username).Exec(); err != nil {
            log.Fatal(err)
            return  err
        }
    }
    return
}

// DeleteUser removes a user, raising ErrDeleteNull if that user was missing.
func (b CassandraAuthBackend) DeleteUser(username string) error {
    session := b.session
    //Probably a better way to do this in one query
    if _, err := b.User(username); err != nil {
        if err = session.Query(`DELETE from users WHERE username = ?`, username).Exec(); err != nil {
            log.Fatal(err)
            return err
        }
    } else {
        return ErrDeleteNull
    }
    return nil
}

// Close cleans up the backend once done with. This should be called before
// program exit.
func (b CassandraAuthBackend) Close() {
	if b.session != nil {
		b.session.Close()
	}
}
// User returns the user with the given username. Error is set to
// ErrMissingUser if user is not found.
/*func (b MongodbAuthBackend) User(username string) (user UserData, e error) {
	var result UserData

	c := b.connect()
	defer c.Database.Session.Close()

	err := c.Find(bson.M{"Username": username}).One(&result)
	if err != nil {
		return result, ErrMissingUser
	}
	return result, nil
}

// Users returns a slice of all users.
func (b MongodbAuthBackend) Users() (us []UserData, e error) {
	c := b.connect()
	defer c.Database.Session.Close()

	err := c.Find(bson.M{}).All(&us)
	if err != nil {
		return us, mkmgoerror(err.Error())
	}
	return
}

// SaveUser adds a new user, replacing if the same username is in use.
func (b MongodbAuthBackend) SaveUser(user UserData) error {
	c := b.connect()
	defer c.Database.Session.Close()

	_, err := c.Upsert(bson.M{"Username": user.Username}, bson.M{"$set": user})
	return err
}

// DeleteUser removes a user. ErrNotFound is returned if the user isn't found.
func (b MongodbAuthBackend) DeleteUser(username string) error {
	c := b.connect()
	defer c.Database.Session.Close()

	// raises error if "username" doesn't exist
	err := c.Remove(bson.M{"Username": username})
	if err == mgo.ErrNotFound {
		return ErrDeleteNull
	}
	return err
}

// Close cleans up the backend once done with. This should be called before
// program exit.
func (b MongodbAuthBackend) Close() {
	if b.session != nil {
		b.session.Close()
	}
}*/
