using Newtonsoft.Json.Linq;
using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using System.Web;
using System.Web.UI.WebControls;

namespace Flipped_Win10
{
    static class URL 
    {
        public const string loginUrl = "http://39.99.190.67:8081/login?";
        public const string registerUrl = "http://39.99.190.67:8081/register?";
        public const string friendRecommendUrl = "http://39.99.190.67:8081/recommendUser";
        public const string friendListUrl = "http://39.99.190.67:8081/friendList";
        public const string addFriendUrl = "http://39.99.190.67:8081/addFriend?";
        public const string deleteFriendUrl = "http://39.99.190.67:8081/deleteFriend?";
    }

    public class friendInfo
    {
        public int Age;
        public string Email;
        public string Hobby;
        public string Photo;
        public string Profession;
        public string RealName;
        public string Region;
        public int UserType;
        public string UserName;
    }

    public static class NetWork
    {
        private static string httpToken = "";
        private static Tuple<string, string> userPwd;

        public static Tuple<string, string> UserPWD
        {
            get{ return userPwd; }
            set { userPwd = value; }
        }

        public static string HttpToken{
            get { return httpToken;}
            set { httpToken = value;}
        }

        public static async Task<Tuple<String, String>> RegisterAysnc(Dictionary<String, String> keyValues) {
            var multForm = new MultipartFormDataContent();
            StringBuilder urlParameters = new StringBuilder();
            string imgPath = keyValues["avatarSource"];

            HttpClient client = new HttpClient();
            foreach (var kv in keyValues)
            {
                if (kv.Key == "avatarSource")
                    continue;
                urlParameters.Append(String.Format("{0}={1}&", kv.Key, kv.Value));
            }
            FileStream fs = File.OpenRead(imgPath);
            multForm.Add(new StreamContent(fs), "photo", Path.GetFileName(imgPath));
            var registerURL = URL.registerUrl + urlParameters.ToString();
            //var response = await client.PostAsync(registerURL, multForm);
            string responseData;
            string statusCode;
            using (var response = await client.PostAsync(registerURL, multForm)) 
            {
                responseData = await response.Content.ReadAsStringAsync();
                statusCode = response.StatusCode.ToString();
            }
            return Tuple.Create<String, String>(responseData, statusCode);
        }

        public static string Login(string username, string pwd) {
            HttpClient client = new HttpClient();
            string urlSuffix = String.Format("username={0}&password={1}", username, pwd);
            var task = client.PostAsync(URL.loginUrl + urlSuffix, null);
            var response = task.Result;
            var content = response.Content.ReadAsStringAsync();
            var result = content.Result;
            return result;
        }

        public static Tuple<string, friendInfo> GetRecommendFriend() 
        {
            if (httpToken == "")
            {
                return null;
            }
            else
            {
                HttpWebRequest request = (HttpWebRequest)WebRequest.Create(URL.friendRecommendUrl);
                request.Method = "GET";
                request.Headers.Add("token", httpToken);
                request.Accept = "*/*";
                request.Proxy = null;
                HttpWebResponse response = (HttpWebResponse)request.GetResponse();
                string res;
                using (StreamReader reader = new StreamReader(response.GetResponseStream(), Encoding.UTF8)) 
                {
                    res = reader.ReadToEnd();
                }
                return ConvertToFriendInfo(res);
            }
        }

        public static Tuple<string, friendInfo> ConvertToFriendInfo(string jsData) 
        {
            JObject data = JObject.Parse(jsData);
            Tuple<string, friendInfo> res = null;
            if (data != null) 
            {
                friendInfo user = new friendInfo
                {
                    Age = (int)data["data"]["Age"],
                    Email = (string)data["data"]["Email"],
                    Photo = (string)data["data"]["Photo"],
                    Profession = (string)data["data"]["Profession"],
                    RealName = (string)data["data"]["RealName"],
                    UserType = (int)data["data"]["UserType"],
                    UserName = (string)data["data"]["Username"],
                    Hobby = (string)data["data"]["Hobby"],
                    Region = (string)data["data"]["Region"]
                };
                string msg = (string)data["message"];
                res = new Tuple<string, friendInfo>(msg, user);
            }
            return res;
        }

        public static IList<string> GetFriendList() 
        {
            HttpWebRequest request = (HttpWebRequest)WebRequest.Create(URL.friendListUrl);
            request.Method = "GET";
            request.Headers.Add("token", httpToken);
            request.Accept = "*/*";
            request.Proxy = null;
            HttpWebResponse response = (HttpWebResponse)request.GetResponse();
            string res;
            using (StreamReader reader = new StreamReader(response.GetResponseStream(), Encoding.UTF8))
            {
                res = reader.ReadToEnd();
            }
            JObject data = JObject.Parse(res);

            Debug.WriteLine(data["data"].GetType().ToString());
            Debug.WriteLine(data["message"]);

            IList<string> friendlist =( data["data"] as JArray ).ToObject<List<string>>();
            Debug.WriteLine(friendlist);

            return friendlist;
        }

        public static void UpdateFriendList(string username)
        {
            IList<string> friendlist = GetFriendList();
            foreach (var friend in friendlist)
            {
                LocalDataBaseOperator.insert(username, friend);
            }
        }

        public static void AddFriend(string friendName)
        {
            HttpWebRequest request = (HttpWebRequest)WebRequest.Create($"{URL.addFriendUrl}friend={friendName}");
            request.Method = "POST";
            request.Headers.Add("token", httpToken);
            request.Accept = "*/*";
            request.Proxy = null;
            HttpWebResponse response = null;
            string res;
            try
            {
                response = (HttpWebResponse)request.GetResponse();
            }
            catch(WebException ex)
            {
                response = ex.Response as HttpWebResponse;
            }
            finally
            {
                using (StreamReader reader = new StreamReader(response.GetResponseStream(), Encoding.UTF8))
                {
                    res = reader.ReadToEnd();
                }
                Debug.WriteLine(res);
                //每次添加完好友就同时更新本地好友数据库列表
                NetWork.UpdateFriendList(NetWork.UserPWD.Item1);
            }

        }

        public static void DeleteFriend(string friendName)
        {
            HttpWebRequest request = (HttpWebRequest)WebRequest.Create($"{URL.deleteFriendUrl}friend={friendName}");
            request.Method = "POST";
            request.Headers.Add("token", httpToken);
            request.Accept = "*/*";
            request.Proxy = null;
            HttpWebResponse response = null;

            try
            {
                response = (HttpWebResponse)request.GetResponse();
            }
            catch (WebException ex)
            {
                response = ex.Response as HttpWebResponse;
            }
            finally
            {
                string res;
                using (StreamReader reader = new StreamReader(response.GetResponseStream(), Encoding.UTF8))
                {
                    res = reader.ReadToEnd();
                }
                Debug.WriteLine(res);
                LocalDataBaseOperator.DeleteFriend(NetWork.UserPWD.Item1, friendName);
            }
        }
    }
}
