   M S B A M L   `   `   `  ����  Flipped-Win10    Flipped_Win10.MainWindow   5   	   W SPresentationCore, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35\ XPresentationFramework, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35R NWindowsBase, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35D 9http://schemas.microsoft.com/winfx/2006/xaml/presentation    5   	   R NSystem.Xaml, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089:x,http://schemas.microsoft.com/winfx/2006/xaml     5   	   @d2http://schemas.microsoft.com/expression/blend/2008     5   	   Jmc;http://schemas.openxmlformats.org/markup-compatibility/2006     5   	   V RMaterialDesignThemes.Wpf, Version=3.1.2.1465, Culture=neutral, PublicKeyToken=nullFmaterialDesign1http://materialdesigninxaml.net/winfx/xaml/themes  5   	   -   5          �Name$  window��5   	    � Title$ 
MainWindow��5   	   $	��521��6   $	��681��6)    � WindowStyle$
 None=�65    � AllowsTransparency .  � WindowStartupLocation$ CenterScreen=�6b    � 
ResizeMode$ NoResize=�6�    � Icon$ win.ico��6�   '  !System.Windows.Shell.WindowChrome   WindowChrome 5	        5
   
      GlassFrameThickness$ 0q�6   	   ResizeBorderThickness$	 10q�6/   
   CaptionHeight$
 50Q�6J    � 	Resources 5         ClosedStoryboard �� 
TargetName W� From	 W� To g� Duration �� TargetProperty g� 	BeginTime	 �� To  ClosedBrush �� 
StartPoint �� EndPoint � Color � Offset��%  &	        &	 �    �� 5   
   $
 grid��6.   .��5      W� $ 1X�6   $ 0X�6'   $ 0:0:0.3D�6.   $9 3OpacityMask.(GradientBrush.GradientStops)[1].Offset�6A   W� 5      $ 1X�6   $ 0X�6'   $ 0:0:0.3D�6.   $ 0:0:0.3X�6A   $9 3OpacityMask.(GradientBrush.GradientStops)[2].Offset�6U   �� 5      $ 	#00000000X�6   $ 0D�6,   $8 2OpacityMask.(GradientBrush.GradientStops)[2].Color�69   5      �� 5   
   $ 0.5,0&�62   $ 0.5,1&�6E   .��5      � $ 	#FF000000��6   $ 0Q�6-   � 5      $ 	#FF000000��6   $ 1Q�6-   � 5      $ 	#FF000000��6   $ 1Q�6-   5      5      ��5      �� 5   
    �� 
TargetType" �6   .R�5      �  �Name$ backGrid��6   .+�5      ��  ��Name$ contentBorder��6   $��10q�60    �� CornerRadius$ 10��6<   ���#fefefc J� Effect 5      3  -System.Windows.Media.Effects.DropShadowEffect  5         
BlurRadius$ 6Q�6+      ShadowDepth$ 0Q�6:      Color$ #B7B7B7��6J    J� OpacityMask 5       !� 5!        !� Visual  5"      �� 5#   "   
���   �! J� SnapsToDevicePixels! .  5%   )   �� *CornerRadius+" �� RelativeSource" �� # �� AncestorType$# BorderL���5&   )   �� *ActualWidth+" �� $# BorderL���5'   )   �� *ActualHeight+" �� $# BorderL�5)      ._�5+      �� .u�5,      �� 5-      5.      5/      $ �� Triggers$ 50      P� 51      % P� Property% �@�IsActive& P� Value& . . �52      �� ' �� 
TargetName$' contentBorder��6   ( �� Property( �@J�Effect) �� Value) 53        54      $ 9Q�6/   $ 0Q�6>   $ #B2B2B2��6N   56      57      59      .��5<      � -   5        $
 grid��5<       5>   
   �� 5?      $ 0.5,1&�6"   $ 0.5,0&�65   .��5@      � $ 	#00000000��6   $ 0Q�61   � 5A      $ 	#00000000��6   $ 1Q�61   � 5B      $ 	#00000000��6   $ 1Q�61   5C      * � Triggers* 5F   
   :� 5G      + :� RoutedEvent$+ Loaded��6   [�5H      �� 5I      .��5J      �� .��5K      W� $ 1X�6.   $ 0X�67   $ 0:0:0.5D�6>   $9 3OpacityMask.(GradientBrush.GradientStops)[1].Offset�6Q   W� 5L      $ 1X�6.   $ 0X�67   $ 0:0:0.5D�6>   $ 0:0:0.2X�6Q   $9 3OpacityMask.(GradientBrush.GradientStops)[2].Offset�6e   �� 5M      $ 	#FF000000X�6-   $ 0D�6<   $8 2OpacityMask.(GradientBrush.GradientStops)[2].Color�6I   5N      5O      5Q      .+�5T   
   �� $
��Left=�6   $	��400��6+   $��/Images/startupImage.png��67   �� 5V   
   $
��Left=�6   $	��100��6+   $��
480,95,0,0q�68   $	��Top=�6L   $	��100��6d   $ ��/Images/default_avatar.png��6p   �� 5X   
   -   5        , ��Name$, 
accountBox��5X      $��Center=�6&   $
��Left=�6H   $��50��6c   $��15(�6o   $��Gaoel)�6}   .  (MaterialDesignThemes.Wpf.TextFieldAssist-   HasClearButton- . )  #MaterialDesignThemes.Wpf.HintAssist.   Hint$. Accout��5Y   G   $��418,240,0,0q�5Z      / �� TextWrapping$
/ Wrap=�6'   $	��Top=�6;   $	��224��6S   $��0,0,0,1q�6_   $��8,0,0,0q�6y   �� 5\   
   $��418,290,0,161q�6   $
��Left=�6,   $��LeftToRight=�6G   0 �� Orientation$0 
Horizontal=�6c   .+�5]      G� -   5        1 G�Name$1 pwdBox��5]      $��Center=�6*   $��50��6M   $��15(�6Y   $��Gaoel)�6g   - . $. Password��5^   F   $	��Top=�5_      $	��177��6)   $��0,0,0,1q�65   $��8,0,0,0q�6O   �� 5a      -   5        2 ��Name$2 loginBtn��5a      $��50��6>   
������   CircleButton#��[ 6b   $��8q�6�   $��0,4,2,0q�6�   5b      � 5d   
   $��418,365,19,0q�6   3 � ColumnDefinitions3 5e      �� 5f      �� 5g      4 � RowDefinitions4 5i      �� 5j      �� 5k      .+�5m      �� -   5        5 ��Name$5 cb_rememberPwd��5m      $��0��6/   $��0��6<   $��记住密码��6L   $��2,0,0,0q�6[   6 d� 	IsChecked$
6 TrueY�6l   $��16(�5n      $��SimSun)�6%   $��Light"�69   $
��Left=�6L   �� 5o      -   5        $5 cb_autoLogin��5o      $��1��6-   $��0��6:   $��自动登录��6J   $��2,0,0,0q�6Y   $
6 TrueY�6j   $��16(�5p      $��SimSun)�6%   $��Light"�69   $
��Left=�6L   �� 5q      -   5        $2 
findPwdBtn��5q      $��1��6)   $��0��69   
���   �$��16(�6Y   $��SimSun)�6g   $��Light"�6{   $��Right=�5r        MaterialDesignFlatButton#��[ 61   $��找回密码��6c   �� 5s      -	   5        $2 registerBtn��5s      $��1��6*   $��1��6:   
���   �$��16(�6Z   $��SimSun)�6h   $��Light"�6|   $��Right=�5t      #��[ 61   $��注册账号��6c   6r   5u      �� 5v   
   -
   5        $
��Left=�5v      $	��Top=�6G     WindowCloseButtonStyle#��[ 6_   $��30��6�   $��60��6�   7   IsHitTestVisibleInChrome7 . $��	602,0,0,0q�6�   5w      5x      