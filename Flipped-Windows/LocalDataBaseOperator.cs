using System;
using System.Collections.Generic;
using System.Data.SQLite;
using System.Diagnostics;
using System.IO;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace Flipped_Win10
{
    public static class LocalDataBaseOperator
    {
        private static string dbPath;
        private static SQLiteConnection Conn = null;

        public static void createDB(string dbName)
        {
            if (string.IsNullOrEmpty(dbPath))
            {
                dbPath = $"./{dbName}";
                string absPath = Path.GetFullPath(dbPath);
                SQLiteConnection.CreateFile(absPath);
                Conn = new SQLiteConnection($"data source = {absPath}");
            }
        }

        public static void addTable(string username)
        {
            if (Conn == null)
            {
                throw new NullReferenceException("haven't connected to SQLite");
            }
            else
            {
                if (Conn.State != System.Data.ConnectionState.Open)
                {
                    Conn.Open();
                }
                SQLiteCommand cmd = new SQLiteCommand
                {
                    Connection = Conn,
                    CommandText = $"create table if not exists {username} (friend varchar);"
                };
                cmd.ExecuteNonQuery();
                Conn.Close();
            }
        }

        public static void insert(string username, string friendName)
        {
            if (Conn == null)
            {
                throw new NullReferenceException("haven't connected to SQLite");
            }
            else
            {
                if (Conn.State != System.Data.ConnectionState.Open)
                {
                    Conn.Open();
                }
                SQLiteCommand cmd = new SQLiteCommand
                {
                    Connection = Conn,
                    CommandText = $"insert or ignore into {username} (friend) values ('{friendName}');"
                };
                int rows = cmd.ExecuteNonQuery();
                Debug.WriteLine($"Succed insert into table {username} with {rows} rows");
                Conn.Close();
            }
        }

        public static IList<string> getFriends(string username)
        {
            IList<string> res = new List<string>();

            if (Conn.State != System.Data.ConnectionState.Open)
            {
                Conn.Open();
            }
            using (SQLiteCommand cmd = new SQLiteCommand
            {
                Connection = Conn,
                CommandText = $"select * from {username};"
            }) {
                SQLiteDataReader reader = cmd.ExecuteReader();
                while (reader.Read())
                {
                    res.Add(reader["friend"] as string);
                }
            }
            return res;
        }

        public static bool IsFriends(string sourceUser, string targetUser)
        {
            IList<string> friendList = getFriends(sourceUser);
            return friendList.Contains(targetUser);
        }

        public static void DeleteFriend(string username, string targetUser)
        {
            if (Conn.State != System.Data.ConnectionState.Open)
            {
                Conn.Open();
            }
            using (SQLiteCommand cmd = new SQLiteCommand
            {
                Connection = Conn,
                CommandText = $"delete from {username} where friend = '{targetUser}';"
            })
            {
                int rows = cmd.ExecuteNonQuery();
                Debug.WriteLine($"Succed delete from table {username} with {rows} rows");
            }
        }
    }
}
