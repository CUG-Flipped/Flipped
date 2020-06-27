using Newtonsoft.Json.Linq;
using System;
using System.Diagnostics;
using System.Windows;
using System.Windows.Input;
using System.Windows.Media;
using System.Windows.Media.Animation;

namespace Flipped_Win10
{
    /// <summary>
    /// MainWindow.xaml 的交互逻辑
    /// </summary>
    public partial class MainWindow : Window
    {
        public MainWindow()
        {
            InitializeComponent();
        }

        private void Window_MouseLeftButtonDown(object sender, MouseButtonEventArgs e)
        {
            this.DragMove();
        }

        private void CloseButton_Click(object sender, RoutedEventArgs e)
        {
            this.IsEnabled = false;
            grid.OpacityMask = this.Resources["ClosedBrush"] as LinearGradientBrush;
            Storyboard std = this.Resources["ClosedStoryboard"] as Storyboard;
            std.Completed += delegate { this.Close(); };
            std.Begin();
        }

        private void RegisterBtn_Click(object sender, RoutedEventArgs e)
        {
            var registerWin = new Register();
            this.Hide();
            registerWin.ShowDialog();
            this.Show();
        }

        private void LoginBtn_Click(object sender, RoutedEventArgs e)
        {
            string username = accountBox.Text;
            string pwd = pwdBox.Password;

            if ( username != "" && pwd != "")
            {
                string result = NetWork.Login(username, pwd);
                JObject jsData = JObject.Parse(result);
                string msg = (string)jsData["message"];
                string tokenStr = (string)jsData["data"]["token"];
                Debug.WriteLine(result);
                NetWork.HttpToken = tokenStr;

                if (msg == "succeed to login") 
                {
                    NetWork.UserPWD = new Tuple<string, string>(username, pwd);
                    var friendRecommentWin = new FriendRecommend();
                    this.Hide();
                    friendRecommentWin.ShowDialog();
                    Environment.Exit(0);
                }
            }
            else 
            {
                MessageBox.Show("请填写完账号和密码后在登录！", "Warning");
            }
        }
    }
}
