using System;
using System.Collections.Generic;
using System.Diagnostics;
using System.IO;
using System.Net;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;
using System.Web.UI.WebControls;

namespace Flipped_Win10
{
    static class URL 
    {
       public static string loginUrl = "http://39.99.190.67:8080/login?";
       public static string registerUrl = "http://39.99.190.67:8080/register?";
    }

    public class NetWork
    {
        private static string httpToken = "";

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
    }
}
